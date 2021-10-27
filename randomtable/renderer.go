package randomtable

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type randomTableRenderer struct {
	nodeRendererFuncsTmp map[ast.NodeKind]renderer.NodeRendererFunc
	tree                 Tree
	namespace            []string
	depth                int
}

// Push a string into the namespace
func (r *randomTableRenderer) Push(name string) {
	r.namespace = append(r.namespace, name)
	r.depth++
}

// Pop off the last name segment from the namespace
func (r *randomTableRenderer) Pop() {
	if r.depth != 0 && len(r.namespace) != 0 {
		r.namespace = r.namespace[:len(r.namespace)-1]
		r.depth--
	}
}

func (r *randomTableRenderer) Name() string {
	return strings.ReplaceAll(strings.ToLower(strings.Join(r.namespace, "/")), " ", "")
}

func NewRandomTableRenderer(tree Tree) renderer.NodeRenderer {
	r := &randomTableRenderer{
		nodeRendererFuncsTmp: map[ast.NodeKind]renderer.NodeRendererFunc{},
		tree:                 tree,
		namespace:            []string{},
		depth:                0,
	}

	return r
}

func (r *randomTableRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindTableHeader, r.renderTableHeader)
	reg.Register(gast.KindTableRow, r.renderTableRow)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(gast.KindTable, r.renderTable)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
}

func (r *randomTableRenderer) renderTableHeader(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		var table Table
		// The first cell cotnains the name of the table
		tablename := string(n.FirstChild().Text(source))
		// Push the header into the namespace when entering the header
		r.Push(tablename)
		// A single header cell is a regular table
		if n.ChildCount() == 1 {
			t := NewRandomTable()
			table = &t

		} else if n.ChildCount() == 2 { // A Rolling table has two cells, name and dice to roll
			diceRoll := string(n.LastChild().Text(source))
			t := NewRollingTable(diceRoll).WithLogger(
				log.WithFields(log.Fields{"table": r.Name()}))
			table = &t
		}

		r.tree.AddTable(r.Name(), table)
	} else {
	}
	return ast.WalkContinue, nil
}
func (r *randomTableRenderer) renderTableRow(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		table, err := r.tree.GetTable(r.Name())
		if err != nil {
			return ast.WalkContinue, fmt.Errorf("Unable to find table: %s", r.Name())
		}
		// A single child row is a regular table
		if n.ChildCount() == 1 {
			table.AddItem(string(n.FirstChild().Text(source)))

		} else if n.ChildCount() == 2 { // The 2nd child is the die result needed for this row
			roll := string(n.LastChild().Text(source))
			singleitem, _ := regexp.MatchString("^[0-9]+$", roll)
			if singleitem {
				i, err := strconv.Atoi(roll)
				if err != nil {
					return ast.WalkContinue, err
				}
				table.AddItem(string(n.FirstChild().Text(source)), i)
				return ast.WalkContinue, nil
			}
			matchRange, _ := regexp.MatchString("^[0-9]+-[0-9]+$", roll)
			if matchRange {
				numRange := strings.Split(roll, "-")
				start, _ := strconv.Atoi(numRange[0])
				end, _ := strconv.Atoi(numRange[1])
				for x := start; x <= end; x++ {
					table.AddItem(string(n.FirstChild().Text(source)), x)
				}
				return ast.WalkContinue, nil
			}

		}

	}
	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderHeading(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Heading)
		headername := string(n.Text(source))
		// Pop until the h
		for n.Level <= r.depth {
			r.Pop()
		}
		r.Push(headername)
	}
	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderTable(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Pop the namespace when we finish with a table
	if !entering {
		r.Pop()
	}
	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderEmphasis(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		t, err := r.tree.GetTable(r.Name())
		if err != nil {
			return ast.WalkContinue, err
		}
		// FIXME: Only hide if this is the header cell
		t.Hidden = true
		r.tree.tables.Put(r.Name(), t)
	}
	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderFencedCodeBlock(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		t := NewTextTable()
		n := node.(*ast.FencedCodeBlock)

		// If the block doesn't have a title, don't treat it as a table since it can't be addressed
		if n.Language(source) == nil {
			return ast.WalkContinue, nil
		}
		title := string(n.Language(source))
		hide := false
		// The table should be marked as hidden
		if strings.HasPrefix(title, "_") && strings.HasSuffix(title, "_") {
			title = strings.TrimPrefix(title, "_")
			title = strings.TrimSuffix(title, "_")
			hide = true
		}
		r.Push(title)
		// Combine all the lines into a single string and use that for the table Item
		var result string
		for _, line := range n.Lines().Sliced(0, n.Lines().Len()) {
			result += string(line.Value(source))
		}
		t.AddItem(result)
		r.tree.AddTable(r.Name(), &t)
		if hide {
			tb, _ := r.tree.GetTable(r.Name())
			tb.Hidden = true
			r.tree.tables.Put(r.Name(), tb)
		}
		r.Pop()
	}

	return ast.WalkContinue, nil
}

package randomtable

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	return strings.ToLower(strings.Join(r.namespace, "/"))
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
}

func (r *randomTableRenderer) renderTableHeader(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		var table Table
		// The first cell cotnains the name of the table
		tablename := string(n.FirstChild().Text(source))
		// A single header cell is a regular table
		if n.ChildCount() == 1 {
			t := NewRandomTable()
			table = &t

		} else if n.ChildCount() == 2 { // A Rolling table has two cells, name and dice to roll
			diceRoll := string(n.LastChild().Text(source))
			t, _ := NewRollingTable(diceRoll)
			table = &t
		}
		// Push the header into the namespace when entering the header
		r.Push(tablename)

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
		t.Hidden = true
		r.tree.tables.Put(r.Name(), t)
	}
	return ast.WalkContinue, nil
}

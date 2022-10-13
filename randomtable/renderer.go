package randomtable

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/justinian/dice"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const ROLL_TABLE_NAME = "ROLLTABLECOLUMN"

type randomTableRenderer struct {
	nodeRendererFuncsTmp map[ast.NodeKind]renderer.NodeRendererFunc
	tree                 Tree
	namespace            []string
	depth                int
	currentTableNames    []string //Names of the tables being rendered
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

func (r *randomTableRenderer) Namespace() string {
	return strings.ReplaceAll(strings.ToLower(strings.Join(r.namespace, "/")), " ", "")
}

// Returns the namespaced name for a given table name
func (r *randomTableRenderer) Name(name string) string {
	if len(r.namespace) == 0 {
		return name
	}
	newName := append(r.namespace, name)
	return strings.ReplaceAll(strings.ToLower(strings.Join(newName, "/")), " ", "")
}

func NewRandomTableRenderer(tree Tree) renderer.NodeRenderer {
	r := &randomTableRenderer{
		nodeRendererFuncsTmp: map[ast.NodeKind]renderer.NodeRendererFunc{},
		tree:                 tree,
		namespace:            []string{},
		depth:                0,
		currentTableNames:    []string{},
	}

	return r
}

func (r *randomTableRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindTableHeader, r.renderTableHeader)
	reg.Register(gast.KindTableRow, r.renderTableRow)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindLink, r.renderLink)
}

func (r *randomTableRenderer) parseHeaderCell(cell ast.Node, col int, source []byte) string {
	text := cell.Text(source)
	diceRoll := ""
	// If we find a dice string, all other columns are for rolling
	if (dice.StdRoller{}.Pattern().Match(text)) {
		diceRoll = string(text)
		r.currentTableNames[col] = ROLL_TABLE_NAME
	} else {
		// Push the header into the namespace when entering the header
		r.currentTableNames[col] = r.Name(string(text))
	}
	return diceRoll
}
func (r *randomTableRenderer) renderTableHeader(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.currentTableNames = make([]string, n.ChildCount())
		childNum := 0
		// Header --Child--> 1st Header Cell --Sibling--> Nth Header Cell
		diceRoll := r.parseHeaderCell(n.FirstChild(), childNum, source)
		sib := n.FirstChild().NextSibling()
		childNum++
		for sib != nil {
			newRoll := r.parseHeaderCell(sib, childNum, source)
			if newRoll != "" {
				diceRoll = newRoll
			}
			sib = sib.NextSibling()
			childNum++
		}
		for _, name := range r.currentTableNames {
			// No table needs to be made for this column
			if name == ROLL_TABLE_NAME {
				continue
			}
			if diceRoll == "" {
				t := NewRandomTable()
				r.tree.AddTable(name, &t, false)
			} else {
				t := NewRollingTable(diceRoll).WithLogger(
					log.WithFields(log.Fields{"table": name}))
				r.tree.AddTable(name, &t, false)
			}
		}
	}
	return ast.WalkContinue, nil
}
func (r *randomTableRenderer) renderTableRow(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		rollColumn := -1
		// Match the row cell to the header cell to determine which table we are referencing
		for x, name := range r.currentTableNames {
			if name == ROLL_TABLE_NAME {
				rollColumn = x
			}
		}

		// Get the text from each cell in the row and hold it in an array
		columns := make([]string, n.ChildCount())
		childCount := 0
		// Row --Child--> 1st Row Cell --Sibling--> Nth Row Cell
		child := n.FirstChild()
		text := string(child.Text(source))
		columns[childCount] = string(text)
		childCount++
		sib := child.NextSibling()
		for sib != nil {
			text := string(sib.Text(source))
			columns[childCount] = text
			childCount++
			sib = sib.NextSibling()
		}
		// Take each column item and add it to the corrisponding table
		// Checks to see if we have dice rolls associated with the table
		for x, text := range columns {
			// we don't need to directly add the roll column to any table
			// just us the value for the other columns
			if x == rollColumn {
				continue
			}
			tableName := r.currentTableNames[x]
			table, err := r.tree.GetTable(tableName)

			if err != nil {
				return ast.WalkContinue, fmt.Errorf("Unable to find table: %s", tableName)
			}
			// Not a rolling table
			if rollColumn == -1 {
				table.AddItem(text)

			} else { // This is a rolling table and we need to use the string from the dice column
				roll := columns[rollColumn]
				// A single number will match this row
				singleitem, _ := regexp.MatchString("^[0-9]+$", roll)
				if singleitem {
					r, err := strconv.Atoi(roll)
					if err != nil {
						return ast.WalkContinue, err
					}
					table.AddItem(text, r)
				}
				// a range of numbers will match this row
				matchRange, _ := regexp.MatchString("^[0-9]+-[0-9]+$", roll)
				if matchRange {
					numRange := strings.Split(roll, "-")
					start, _ := strconv.Atoi(numRange[0])
					end, _ := strconv.Atoi(numRange[1])
					for r := start; r <= end; r++ {
						table.AddItem(text, r)
					}
				}
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

func (r *randomTableRenderer) renderEmphasis(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// emphasis -> cell -> Header == this is a table name
		switch node.Parent().Parent().(type) {
		case *gast.TableHeader, *ast.FencedCodeBlock:
			break
		default:
			return ast.WalkContinue, nil
		}
		name := string(node.Text(source))
		name = r.Name(name)
		t, err := r.tree.GetTable(name)
		if err != nil {
			return ast.WalkContinue, err
		}
		t.Hidden = true
		r.tree.tables.Put(name, t)
	}
	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderFencedCodeBlock(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.FencedCodeBlock)

		// If the block doesn't have a title, don't treat it as a table since it can't be addressed
		if n.Language(source) == nil {
			return ast.WalkContinue, nil
		}
		t := NewTextTable()
		title := string(n.Language(source))
		hidden := false
		// The table should be marked as hidden
		if strings.HasPrefix(title, "_") && strings.HasSuffix(title, "_") {
			title = strings.TrimPrefix(title, "_")
			title = strings.TrimSuffix(title, "_")
			hidden = true
		}
		title = r.Name(title)
		// Combine all the lines into a single string and use that for the table Item
		var result string
		for _, line := range n.Lines().Sliced(0, n.Lines().Len()) {
			result += string(line.Value(source))
		}
		t.AddItem(result)
		r.tree.AddTable(title, &t, hidden)
	}

	return ast.WalkContinue, nil
}

func (r *randomTableRenderer) renderLink(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// link -> Header == this is a valid alias
		// switch node.Parent().(type) {
		// case *gast.TableHeader, *ast.FencedCodeBlock:
		// 	break
		// default:
		// 	return ast.WalkContinue, nil
		// }
		n := node.(*ast.Link)
		label := string(n.Text(source))
		url := string(n.Destination)
		r.tree.AddLink(r.Name(label), url)
	}
	return ast.WalkContinue, nil
}

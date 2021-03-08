package randomtable

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/dghubble/trie"
	"github.com/justinian/dice"
)

// Tree holds lookup tables and allows for retrieveing tables as well as items from tables.
// In addition, the tree handles the rendering of any templates that are a part of a table
type Tree struct {
	tables         *trie.PathTrie
	lookupDepth    int
	maxLookupDepth int
}

// TableNode embeds the table that was created and adds meta-data for use in the tree
type TableNode struct {
	Table
	Hidden bool
}

func NewTree() Tree {
	return Tree{
		tables:         trie.NewPathTrie(),
		lookupDepth:    0,
		maxLookupDepth: 100,
	}
}

// AddTable adds the given table with the given name
func (t *Tree) AddTable(name string, table Table) {
	name = strings.ToLower(name)
	t.tables.Put(name, TableNode{Table: table})
}

// GetTable returns the table with the given name in the tree
func (t *Tree) GetTable(name string) (TableNode, error) {
	name = strings.ToLower(name)
	table := t.tables.Get(name)
	//fmt.Printf("Table: %s", table)
	if table == nil {
		return TableNode{}, fmt.Errorf("%s table not found", name)
	}
	tb := table.(TableNode)
	switch tableTyped := tb.Table.(type) {
	default:
		tb.Table = tableTyped
		return tb, nil
	}
}

// ListTables will return a list of all the tables that aren't marked as hiddn, that are loaded in the tree, and that start with the given prefix
func (t *Tree) ListTables(prefix string) []string {
	tables := []string{}
	t.tables.Walk(func(key string, value interface{}) error {
		tb, ok := value.(TableNode)
		if ok && strings.HasPrefix(key, prefix) && !tb.Hidden {
			tables = append(tables, key)
		}
		return nil
	})
	return tables
}

// GetItem retreieves an item from a table and will render any items
// that include templates. Wraps getItem for loop detection
func (t *Tree) GetItem(table string) (string, error) {

	item, err := t.getItem(table)
	//reset the lookup now that we've finished
	t.lookupDepth = 0
	return item, err
}

// renderItem will render any templates for a given item. Table is the path the item was
// found on to allow for lookups using relative paths
func (t *Tree) renderItem(item string, table string) (string, error) {
	funcMap := template.FuncMap{
		"lookup": t.getLookup(table),
		"roll":   t.roll,
	}
	tmpl, err := template.New("item").Funcs(funcMap).Parse(item)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, nil)
	return buf.String(), nil

}

// getLookup provides a function for retrieving items from other tables.
// It uses a closure to provide the calling table to allow relative pathing
func (t *Tree) getLookup(callingTable string) func(string, ...interface{}) string {
	return func(item string, rolls ...interface{}) string {
		// number of times to roll
		var times int

		// replace the relative path with the full path
		if strings.HasPrefix(item, "./") {
			tablePaths := strings.Split(callingTable, "/")
			tablePrefix := strings.Join(tablePaths[0:len(tablePaths)-1], "/")
			item = strings.Replace(item, "./", tablePrefix+"/", 1)
		}
		if len(rolls) == 0 {
			times = 1
		} else {
			var err error
			switch r := rolls[0].(type) {
			case string:
				times, err = strconv.Atoi(r)
				if err != nil {
					times = 1
				}
			case int:
				times = r
			}
		}

		// checking for a loop
		if t.lookupDepth >= t.maxLookupDepth {
			return item
		}
		t.lookupDepth++

		result := []string{}
		for x := 1; x <= times; x++ {
			i, err := t.getItem(item)
			if err == nil {
				result = append(result, i)
			}
		}
		return strings.Join(result, ", ")
	}

}

// roll is a template function for rolling dice on a table
func (t *Tree) roll(d string) string {
	result, _, err := dice.Roll(d)
	if err != nil {
		return d
	}
	return strconv.Itoa(result.Int())
}

func (t *Tree) getItem(table string) (string, error) {
	tb, err := t.GetTable(table)
	if err != nil {
		return "", err
	}
	item := tb.GetItem()
	return t.renderItem(item, table)
}

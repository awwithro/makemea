package randomtable

import (
	"bytes"
	"errors"
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
	t.tables.Put(name, table)
}

// GetTable returns the table with the given name in the tree
func (t *Tree) GetTable(name string) (Table, error) {
	name = strings.ToLower(name)
	table := t.tables.Get(name)
	if table == nil {
		return nil, fmt.Errorf("%s table not found", name)
	}
	switch tableTyped := table.(type) {
	case *RandomTable:
		return tableTyped, nil
	case *RollingTable:
		return tableTyped, nil
	}
	return nil, errors.New("Unable to determine table type")
}

// ListTables will return a list of all the tables that are loaded in the tree that start with the given prefix
func (t *Tree) ListTables(prefix string) []string {
	tables := []string{}
	t.tables.Walk(func(key string, value interface{}) error {
		_, ok := value.(Table)
		if ok && strings.HasPrefix(key, prefix) {
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
func (t *Tree) getLookup(table string) func(string) string {
	return func(item string) string {
		// replace the relative path with the full path
		if strings.HasPrefix(item, "./") {
			tablePaths := strings.Split(table, "/")
			pathToTable := strings.Join(tablePaths[0:len(tablePaths)-1], "/")
			item = strings.Replace(item, "./", pathToTable+"/", 1)
		}
		// checking for a loop
		if t.lookupDepth >= t.maxLookupDepth {
			return item
		}
		t.lookupDepth++
		i, _ := t.getItem(item)
		return i
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

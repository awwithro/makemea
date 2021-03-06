package randomtable

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/dghubble/trie"
)

type RandomTableTree struct {
	tables         *trie.PathTrie
	lookupDepth    int
	maxLookupDepth int
}

func NewRandomTableTree() RandomTableTree {
	return RandomTableTree{
		tables:         trie.NewPathTrie(),
		lookupDepth:    0,
		maxLookupDepth: 100,
	}
}

func (t *RandomTableTree) AddTable(name string, table Table) {
	name = strings.ToLower(name)
	t.tables.Put(name, table)
}

func (t *RandomTableTree) GetTable(name string) (Table, error) {
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

func (t *RandomTableTree) ListTables(prefix string) []string {
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
func (t *RandomTableTree) GetItem(table string) (string, error) {

	item, err := t.getItem(table)
	//reset the lookup now that we've finished
	t.lookupDepth = 0
	return item, err
}

func (t *RandomTableTree) renderItem(item string) (string, error) {
	funcMap := template.FuncMap{
		"lookup": t.lookup,
	}
	tmpl, err := template.New("item").Funcs(funcMap).Parse(item)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, nil)
	return buf.String(), nil

}

// lookup is a template function for looking up items on other tables
func (t *RandomTableTree) lookup(item string) string {
	// checking for a loop
	if t.lookupDepth >= t.maxLookupDepth {
		return item
	}
	t.lookupDepth++
	i, _ := t.getItem(item)
	return i
}

func (t *RandomTableTree) getItem(table string) (string, error) {
	tb, err := t.GetTable(table)
	if err != nil {
		return "", err
	}
	item := tb.GetItem()
	return t.renderItem(item)
}

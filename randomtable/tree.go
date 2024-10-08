package randomtable

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/dghubble/trie"
	"github.com/justinian/dice"
	log "github.com/sirupsen/logrus"
)

// Tree holds lookup tables and allows for retrieveing tables as well as items from tables.
// In addition, the tree handles the rendering of any templates that are a part of a table
type Tree struct {
	tables         *trie.PathTrie
	maxLookupDepth int
	formatter      Formatter
}

// TableNode embeds the table that was created and adds meta-data for use in the tree
type TableNode struct {
	Table
	Hidden bool
}

// A link to another table
type LinkNode struct {
	Link string
}

func NewTree() Tree {
	return Tree{
		tables:         trie.NewPathTrie(),
		maxLookupDepth: 100,
		formatter:      StringFormatter{},
	}
}

// AddTable adds the given table with the given name.
// Names have spaces removed and turned to lowercase
func (t *Tree) AddTable(name string, table Table, hidden bool) {
	name = strings.ReplaceAll(strings.ToLower(name), " ", "")

	// check for existing table
	_, err := t.GetItem(name)
	// no err means we got a table
	if err == nil {
		log.WithField("table", name).Warn("Duplicate table entered")
	}

	t.tables.Put(name, TableNode{Table: table, Hidden: hidden})
}

// AddLink adds a reference to another table
func (t *Tree) AddLink(name, table string) {
	name = strings.ReplaceAll(strings.ToLower(name), " ", "")

	// check for existing table
	_, err := t.GetItem(name)
	// no err means we got a table
	if err == nil {
		log.WithField("table", name).Warn("Duplicate table entered")
	}

	t.tables.Put(name, LinkNode{Link: table})
}

// GetTable returns the table with the given name in the tree
func (t *Tree) GetTable(name string) (TableNode, string, error) {
	name = strings.ReplaceAll(strings.ToLower(name), " ", "")
	table := t.tables.Get(name)
	if table == nil {
		return TableNode{},"", fmt.Errorf("%s table not found", name)
	}
	switch tb := table.(type) {
	case TableNode:
		switch tableTyped := tb.Table.(type) {
		default:
			tb.Table = tableTyped
			return tb, name, nil
		}
	case LinkNode:
		linkedTable, name, err := t.GetTable(tb.Link)
		return linkedTable, name, err
	default:
		return TableNode{}, "",fmt.Errorf("unknown Table Node: %v", tb)
	}
}

// ListTables will return a sorted list of all the tables that are loaded in the tree, and that start with the given prefix.
// showHidden is used to toggle weather to show tables that are marked as hidden
func (t *Tree) ListTables(prefix string, showHidden bool) []string {
	tables := sort.StringSlice{}
	t.tables.Walk(func(key string, value interface{}) error {
		switch tb := value.(type) {
		case TableNode:
			if strings.HasPrefix(key, prefix) && (showHidden || !tb.Hidden) {
				tables = append(tables, key)
			}
		case LinkNode:
			if strings.HasPrefix(key, prefix) {
				tables = append(tables, key)
			}
		}

		return nil
	})
	tables.Sort()
	return tables
}

// GetItem retrieves an item from a table and will render any items
// that include templates.
func (t *Tree) GetItem(table string) (string, error) {
	tb, name, err := t.GetTable(table)
	if err != nil {
		return "", err
	}
	item := tb.GetItem()
	item = t.formatter.Format(item, name)
	return t.renderItem(item, name)
}

// renderItem will render any templates for a given item. Table is the path the item was
// found on to allow for lookups using relative paths
func (t *Tree) renderItem(item string, table string) (string, error) {
	funcMap := template.FuncMap{
		"lookup": t.getLookup(table),
		"roll":   t.roll,
		"fudge":  t.getFudge(table),
		"pick": pickItem,
		"chance": chance,
	}
	mergedFuncMaps := sprig.FuncMap()
	for k, v := range funcMap {
		mergedFuncMaps[k] = v
	}
	tmpl, err := template.New("item").Funcs(template.FuncMap(mergedFuncMaps)).Parse(item)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func pickItem(items ...string) string {
	return items[rand.Intn(len(items))]
}

func chance(chance float32, fallback, original string) string{
	if rand.Float32() <= chance{
		return original
	}
	return fallback
}

// getLookup provides a function for retrieving items from other tables.
// It uses a closure to provide the calling table to allow relative pathing
func (t *Tree) getLookup(callingTable string) func(string, ...interface{}) (string, error) {
	return func(item string, rolls ...interface{}) (string, error) {
		item = resolvePaths(callingTable, item)
		// number of times to roll
		times := parseRollCount(rolls)
		result := []string{}
		for x := 1; x <= times; x++ {
			i, err := t.GetItem(item)
			if err != nil {
				return "", err
			}
			result = append(result, i)
		}
		return strings.Join(result, ", "), nil
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

func (t *Tree) ValidateTables() {
	t.tables.Walk(func(key string, value interface{}) error {
		// Call each table to validate itself
		if tb, ok := value.(TableNode); ok {
			tb.Validate()

			// Get all the items and check that they are valid.
			items := tb.AllItems()
			for _, item := range items {
				_, err := t.renderItem(item, key)
				if err != nil {
					log.WithField("table", key).Warn(err)
				}
			}
		}
		return nil
	})
}

// fudge performs a lookup on the given table but uses and alternate dice string
func (t *Tree) getFudge(callingTable string) func(string, string, ...interface{}) (string, error) {
	return func(table, dicestr string, rolls ...interface{}) (string, error) {
		table = resolvePaths(callingTable, table)
		tb, _,err := t.GetTable(table)
		if err != nil {
			return "", err
		}
		times := parseRollCount(rolls)
		var newTable = NewRollingTable(dicestr)
		switch rt := tb.Table.(type) {
		case *RollingTable:
			for k, v := range rt.items {
				newTable.items[k] = v
			}
			newTable.dicestr = dicestr
		// Wonky as items is two different types in these tables
		case *RandomTable:
			for k, v := range rt.items {
				newTable.items[k+1] = v
			}
			newTable.dicestr = dicestr
		}

		result := []string{}
		for x := 1; x <= times; x++ {
			i := newTable.GetItem()
			item, _ := t.renderItem(i, table)
			result = append(result, t.formatter.Format(item, table))
		}
		return strings.Join(result, ", "), nil
	}
}

func resolvePaths(callingTable, table string) string {
	// replace the relative path with the full path
	if strings.HasPrefix(table, "./") {
		tablePaths := strings.Split(callingTable, "/")
		tablePrefix := strings.Join(tablePaths[0:len(tablePaths)-1], "/")
		table = strings.Replace(table, "./", tablePrefix+"/", 1)
	}
	return table
}

func parseRollCount(rolls []interface{}) int {
	var times int
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
	return times
}

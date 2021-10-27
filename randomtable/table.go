package randomtable

import (
	"math/rand"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Table interface {
	GetItem() string
	AddItem(string, ...int)
	Validate()
	AllItems() []string
	GetTable(*tablewriter.Table, string) *tablewriter.Table
}

type RandomTable struct {
	items []string
	seed  int
}

func (r *RandomTable) GetItem() string {
	randomIndex := rand.Intn(len(r.items))
	return r.items[randomIndex]
}

func (r *RandomTable) AddItem(item string, n ...int) {
	r.items = append(r.items, item)
}

func (r *RandomTable) Validate() {
	return
}

func (r RandomTable) AllItems() []string {
	return r.items
}

func (r RandomTable) GetTable(t *tablewriter.Table, name string) *tablewriter.Table {
	for _, item := range r.items {
		t.Append([]string{item})
	}
	t.SetHeader([]string{name})
	return t
}

func NewRandomTable() RandomTable {
	t := RandomTable{
		items: []string{},
		seed:  time.Now().Nanosecond(),
	}
	rand.Seed(int64(t.seed))
	return t
}

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
	rand *rand.Rand
}

func (r *RandomTable) GetItem() string {
	randomIndex := r.rand.Intn(len(r.items))
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
	seed := time.Now().Nanosecond()
	t := RandomTable{
		items: []string{},
		seed:  seed,
		rand: rand.New(rand.NewSource(int64(seed))),
	}
	
	return t
}

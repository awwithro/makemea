package randomtable

import (
	"math/rand"
	"time"
)

type Table interface {
	GetItem() string
	AddItem(string, ...int)
	Validate()
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

func NewRandomTable() RandomTable {

	t := RandomTable{
		items: []string{},
		seed:  time.Now().Nanosecond(),
	}
	rand.Seed(int64(t.seed))
	return t
}

package randomtable

import (
	"math/rand"
	"time"
)

type Table interface {
	GetItem() string
	AddItem(string, ...int)
}

type RandomTable struct {
	items []string
	seed  int
}

func (r *RandomTable) GetItem() string {
	rand.Seed(int64(r.seed))
	randomIndex := rand.Intn(len(r.items))
	return r.items[randomIndex]
}

func (r *RandomTable) AddItem(item string, n ...int) {
	r.items = append(r.items, item)
}

func NewRandomTable() RandomTable {
	return RandomTable{
		items: []string{},
		seed:  time.Now().Nanosecond(),
	}
}

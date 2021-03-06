package randomtable

import (
	"math/rand"
	"time"

	"github.com/justinian/dice"
)

type RollingTable struct {
	items   map[int]string
	dicestr string
	seed    int
}

func (r *RollingTable) GetItem() string {
	rand.Seed(int64(r.seed))
	result, _, _ := dice.Roll(r.dicestr)
	return r.items[result.Int()]
}

func (r *RollingTable) AddItem(item string, pos ...int) {
	for _, x := range pos {
		r.items[x] = item
	}

}

func NewRollingTable(d string) (RollingTable, error) {

	table := RollingTable{
		items:   map[int]string{},
		dicestr: d,
		seed:    time.Now().Nanosecond(),
	}
	return table, nil
}

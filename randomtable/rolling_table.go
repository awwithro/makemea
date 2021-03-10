package randomtable

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/justinian/dice"
	log "github.com/sirupsen/logrus"
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
		if _, exists := r.items[x]; exists {
			log.Warnf("Duplicate item: %s for roll %v", item, x)
		}
		r.items[x] = item
	}

}

// Validate that all numbers in the table are represented, that all numbers can be rolled, and there are no overlapping rolls
func (r *RollingTable) Validate() *multierror.Error {
	result := &multierror.Error{}
	count, sides, err := parseDiceString(r.dicestr)
	if err != nil {
		multierror.Append(result, err)
	}
	minRoll := count
	maxRoll := count * sides
	keys := []int{}

	//look for rolls that can't be reached
	for k, _ := range r.items {
		if k < minRoll || k > maxRoll {
			result = multierror.Append(result, fmt.Errorf("%v is outside of the dice range", k))
		}
		keys = append(keys, k)
	}

	// Determine all possible numbers that can be rolled
	var allRolls = make([]int, maxRoll-minRoll+1)
	for i, x := 0, minRoll; x <= maxRoll; i, x = i+1, x+1 {
		allRolls[i] = x
	}

	// Look for rolls that can't be made. Table is missing numbers
	diff := difference(allRolls, keys)
	for _, roll := range diff {
		result = multierror.Append(result, fmt.Errorf("%v is not rollable", roll))
	}

	return result
}

func NewRollingTable(d string) (RollingTable, error) {
	table := RollingTable{
		items:   map[int]string{},
		dicestr: d,
		seed:    time.Now().Nanosecond(),
	}
	return table, nil
}

func parseDiceString(dicestr string) (int, int, error) {
	pattern := dice.StdRoller{}.Pattern()
	matches := pattern.FindStringSubmatch(dicestr)

	count, err := strconv.ParseInt(matches[1], 10, 0)
	if err != nil {
		return 0, 0, err
	}

	sides, err := strconv.ParseInt(matches[2], 10, 0)
	if err != nil {
		return 0, 0, err
	}

	return int(count), int(sides), nil
}

func difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

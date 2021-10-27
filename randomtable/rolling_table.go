package randomtable

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/awwithro/makemea/util"
	"github.com/justinian/dice"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type RollingTable struct {
	items   map[int]string
	dicestr string
	seed    int
	log     log.Entry
}

func (r *RollingTable) GetItem() string {
	result, _, _ := dice.Roll(r.dicestr)
	return r.items[result.Int()]
}

func (r *RollingTable) AddItem(item string, pos ...int) {
	for _, x := range pos {
		if _, exists := r.items[x]; exists {
			r.log.Warnf("Duplicate item: %s for roll %v", item, x)
		}
		r.items[x] = item
	}
}

func (r RollingTable) AllItems() []string {
	values := []string{}
	for _, val := range r.items {
		values = append(values, val)
	}
	return util.DeDupe(values)
}

// Validate that all numbers in the table are represented, that all numbers can be rolled, and there are no overlapping rolls
func (r *RollingTable) Validate() {
	count, sides, err := parseDiceString(r.dicestr)
	if err != nil {
		r.log.Warn(err)
	}
	minRoll := count
	maxRoll := count * sides
	keys := []int{}

	//look for rolls that can't be reached
	for k := range r.items {
		if k < minRoll || k > maxRoll {
			r.log.Warnf("%v is outside of the dice range", k)
		}
		keys = append(keys, k)
	}

	// Determine all possible numbers that can be rolled
	var allRolls = make([]int, maxRoll-minRoll+1)
	for i, x := 0, minRoll; x <= maxRoll; i, x = i+1, x+1 {
		allRolls[i] = x
	}

	// Look for rolls that can't be made. Table is missing numbers
	diff := util.Difference(allRolls, keys)
	for _, roll := range diff {
		r.log.Warnf("%v is not rollable", roll)
	}
}

func (r RollingTable) GetTable(t *tablewriter.Table, name string) *tablewriter.Table {
	for roll, item := range r.items {
		t.Append([]string{item, strconv.Itoa(roll)})
	}
	t.SetHeader([]string{name, r.dicestr})
	return t
}

func NewRollingTable(d string) RollingTable {
	table := RollingTable{
		items:   map[int]string{},
		dicestr: d,
		seed:    time.Now().Nanosecond(),
		log:     *log.NewEntry(log.StandardLogger()),
	}
	rand.Seed(int64(table.seed))
	return table
}

func (r RollingTable) WithLogger(logger *log.Entry) RollingTable {
	r.log = *logger
	return r
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

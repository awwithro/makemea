package randomtable

import (
	"testing"
)

func TestRollingTable(t *testing.T) {
	r := NewRollingTable("2d4")

	r.AddItem("Hello", 2, 3, 4, 5, 6, 7, 8)
	if r.GetItem() != "Hello" {
		t.Error("Didn't get expected item")
	}
}

type DiceResult struct {
	count int
	sides int
	err   error
}
type ParseDiceStringTest struct {
	dicestr  string
	expected DiceResult
}

func TestParseDiceString(t *testing.T) {
	cases := []ParseDiceStringTest{
		{
			dicestr:  "1d6",
			expected: DiceResult{count: 1, sides: 6, err: nil},
		},
		{
			dicestr:  "2d10",
			expected: DiceResult{count: 2, sides: 10, err: nil},
		},
	}
	for _, c := range cases {
		actualCount, actualSides, err := parseDiceString(c.dicestr)
		if actualCount != c.expected.count || actualSides != c.expected.sides || err != c.expected.err {
			t.Errorf("Expected: %v, %v, %v but got %v, %v, %v", c.expected.count, c.expected.sides, c.expected.err, actualCount, actualSides, err)
		}
	}
}

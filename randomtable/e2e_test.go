package randomtable

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TestCases struct {
	table     string
	tablePath string
	name      string
	expected  []string
}

func TestHeaderLookups(t *testing.T) {
	tests := []TestCases{
		{table: `
| Color  |
| ------ |
| Blue   |
| Red    |
| Yellow |
`,
			tablePath: "color",
			name:      "Test Simple Lookup",
			expected: []string{
				"Blue", "Red", "Yellow",
			},
		},
		{
			table: `
# Places

| Country |
| ------- |
| USA     |
| Mexico  |
| Canada  |
`,
			tablePath: "places/country",
			name:      "Test nested table",
			expected: []string{
				"USA", "Mexico", "Canada",
			},
		},
		{
			table: `
# People
# Places
## Castle

| Name                 |
| -------------------- |
| Roogna               |
| Grayskull            |
| Castle AARrrrrgghhhh |
| Edinburgh            |
| Neuschwanstein       |
`,
			tablePath: "places/castle/name",
			name:      "Test multiple nestings",
			expected: []string{
				"Roogna", "Grayskull", "Castle AARrrrrgghhhh", "Edinburgh", "Neuschwanstein",
			},
		},
		{
			table: `
# Things

| Item   | 2d4 |
| ------ | --- |
| Dagger | 2   |
| Coin   | 3-6 |
| Gem    | 7   |
| Sword  | 8   |
`,
			tablePath: "things/item",
			name:      "Test Rolling",
			expected: []string{
				"Dagger", "Coin", "Gem", "Sword",
			},
		},
		{
			table: `
# Things

| Item   | 2d4 |
| ------ | --- |
| Dagger | 2   |
| Coin   | 3-6 |
| Gem    | 7   |
| Sword  | 8   |

| Fancy                          |
| ------------------------------ |
| Shiny {{lookup "things/item"}} |
`,
			tablePath: "things/fancy",
			name:      "Test lookup of other tables",
			expected: []string{
				"Shiny Dagger", "Shiny Coin", "Shiny Gem", "Shiny Sword",
			},
		},
		{
			table: `
# Nested

| Lookup               |
| -------------------- |
| {{lookup "./table"}} |

| Table                        |
| ---------------------------- |
| Foo                          |
| {{lookup "./subnest/table"}} |

## Subnest
| Table                    |
| ------------------------ |
| Bar                      |
| {{lookup "./subtable" }} |

| Subtable |
| -------- |
| Baz      |
`,
			name:      "Test relative pathing",
			tablePath: "nested/lookup",
			expected: []string{
				"Foo", "Bar", "Baz",
			},
		},
		{
			// No backticks in string literals :(
			table:     fmt.Sprint("\n``` test\ntest\n```\n\n"),
			name:      "Test text table",
			tablePath: "test",
			expected: []string{
				"test\n",
			},
		},
		{
			table: `
| t1  |
| --- |
| one |

| t2                |
| ----------------- |
| {{lookup "t1" 2}} |
			`,
			name:      "Test lookup with counts",
			tablePath: "t2",
			expected: []string{
				"one, one",
			},
		},
		{
			table: `
| t1  |
| --- |
| one |

| t2                  |
| ------------------- |
| {{lookup "t1" "2"}} |
			`,
			name:      "Test lookup with counts as strings",
			tablePath: "t2",
			expected: []string{
				"one, one",
			},
		},
		{
			table: `	
| t1    | 1d4 |
| ----- | --- |
| one   | 1   |
| two   | 2   |
| three | 3   |
| four  | 4   |
| five  | 5   |
| six   | 6   |


| t2                   |
| -------------------- |
| {{fudge "t1" "4d1"}} |
			`,
			name:      "Test fudge works on a roll table",
			tablePath: "t2",
			expected: []string{
				"four",
			},
		},
		{
			table: `	
| t1    | 1d4 |
| ----- | --- |
| one   | 1   |
| two   | 2   |
| three | 3   |
| four  | 4   |
| five  | 5   |
| six   | 6   |


| t2                   |
| -------------------- |
| {{fudge "t1" "4d1" 2}} |
			`,
			name:      "Test fudge works on a roll table with a count",
			tablePath: "t2",
			expected: []string{
				"four, four",
			},
		},
	}
	for _, tc := range tests {
		tree := NewTree()
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(NewRandomTableRenderer(tree), 1))),
		)
		var buf bytes.Buffer
		if err := md.Convert(bytes.NewBufferString(tc.table).Bytes(), &buf); err != nil {
			t.Error(err)
		}
		actual, err := tree.GetItem(tc.tablePath)
		if err != nil {
			t.Errorf("Test: %v, Error: %v", tc.name, err)
		}
		found := false
		for _, exepctedItem := range tc.expected {
			if actual == exepctedItem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s Expected to find one of %s but got %s", tc.name, tc.expected, actual)
		}
	}

}

func TestListTables(t *testing.T) {
	tests := []TestCases{
		{
			table:     listTest,
			tablePath: "",
			name:      "Test Listing with nested and hidden tables",
			expected: []string{
				"t1", "h1/t2", "h1/t3", "h1/h2/t4",
			},
		},
		{
			table:     listTest,
			tablePath: "h1",
			name:      "Test prefix based listing",
			expected: []string{
				"h1/t2", "h1/t3", "h1/h2/t4",
			},
		},
	}
	for _, tc := range tests {
		tree := NewTree()
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(NewRandomTableRenderer(tree), 1))),
		)
		var buf bytes.Buffer
		if err := md.Convert(bytes.NewBufferString(tc.table).Bytes(), &buf); err != nil {
			t.Error(err)
		}
		actual := tree.ListTables(tc.tablePath, false)
		expectedSorted := sort.StringSlice(tc.expected)
		expectedSorted.Sort()
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("%s: Expected to find %s but got %s", tc.name, expectedSorted, actual)
		}
	}
}

const listTest = `
| t1  |
| --- |

| _t5_ |
| ---- |

# h1
| t2  |
| --- |

| t3  |
| --- |

## h2

| t4  |
| --- |

`

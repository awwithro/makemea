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
# Things

| 1d4 | fancy  |
| --- | ---    |
| 1   | Dagger |
| 2   | Coin   |
| 3   | Gem    |
| 4   | Sword  |
`,
			tablePath: "things/fancy",
			name:      "roll column can be anywhere",
			expected: []string{
				"Dagger", "Coin", "Gem", "Sword",
			},
		},
		{
			table: `

| 1d4 | fancy  | another|
| --- | ---    | ---    |
| 1   | Dagger | Dagger |
| 2   | Coin   | Coin   |
| 3   | Gem    | Gem    |
| 4   | Sword  | Sword  |
`,
			tablePath: "another",
			name:      "more than one table with a roll column can be used",
			expected: []string{
				"Dagger", "Coin", "Gem", "Sword",
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
| t1    |
| ----- |
| one   |
| two   |
| three |
| four  |
| five  |
| six   |


| t2                   |
| -------------------- |
| {{fudge "t1" "4d1"}} |
			`,
			name:      "Test fudge works on a non-rolling table",
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
		{
			table:     multiTable,
			name:      "Test we can have two tables in one",
			tablePath: "t1",
			expected: []string{
				"one",
			},
		},
		{
			table:     multiTable,
			name:      "Test we can have two tables in one",
			tablePath: "t2",
			expected: []string{
				"two",
			},
		},
		{
			table: `
# Nested

[link](nested/subnest/table)

## Subnest

| Table |
| ------|
| Bar   |
`,
			name:      "Test links to other tables work",
			tablePath: "nested/link",
			expected: []string{
				"Bar",
			},
		},
		{
			table: `
| t1 |
| --- |
| test{{pick "ing" "s" "ed"}} |
`,
			name: "Test pick item template",
			tablePath: "t1",
			expected: []string{
				"testing","tests","tested",
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
			t.Errorf("Test: %v, Error: %v. Found: %v", tc.name, err, tree.ListTables("", true))
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
		{
			table:     multiTable,
			tablePath: "",
			name:      "Test multiple rollTables in a markdown table",
			expected: []string{
				"t1", "t2",
			},
		},
		{
			table:     multiRoller,
			tablePath: "",
			name:      "Test multiple rollTables in a table w/ dice column",
			expected: []string{
				"t1", "t2",
			},
		},
		{
			table:     hiddenText,
			tablePath: "",
			name:      "Test text blocks can be hidden",
			expected:  []string{},
		},
		{
			table:     linkTest,
			tablePath: "",
			name:      "Test links are listed",
			expected: []string{
				"t1", "t2",
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

func TestHtmlFormattedTables(t *testing.T) {
	tests := []TestCases{
		{
			table:     multiTable,
			tablePath: "t1",
			name:      "Test Html gets formatted properly",
			expected: []string{
				"<randomElement table='t1'>one</randomElement>",
			},
		},
		{
			table:     nestedTable,
			tablePath: "t1",
			name:      "Test Nested Html gets formatted properly",
			expected: []string{
				"<randomElement table='t1'>one: <randomElement table='t2'>two</randomElement></randomElement>",
			},
		},
	}
	for _, tc := range tests {
		tree := NewTree().WithHtmlFormatter()
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
		actual, _ := tree.GetItem(tc.tablePath)
		if !reflect.DeepEqual([]string{actual}, tc.expected) {
			t.Errorf("%s: Expected to find %s but got %s.", tc.name, tc.expected, actual)
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
const multiTable = `
|t1 |t2 |
|---|---|
|one|two|
`
const multiRoller = `
|t1 |t2 |1d1|
|---|---|---|
|one|two|1  |
`
const hiddenText = "``` _test_\n```"

const nestedTable = `
| t1  |
| --- |
| one: {{lookup "t2" }}|

| t2  |
| --- |
| two |
`
const linkTest = `
[t1](t2)
[t3](http://url.com)
[t3](https://url.com)
|t2|
|---|
`

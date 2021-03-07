package randomtable

import (
	"bytes"
	"reflect"
	"sort"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type TestCases struct {
	tablePath string
	expected  []string
}

func TestHeaderLookups(t *testing.T) {
	tree := NewTree()
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewRandomTableRenderer(tree), 1))),
	)
	if err := md.Convert(bytes.NewBufferString(testmakrkwon).Bytes(), &buf); err != nil {
		t.Error(err)
	}
	tests := []TestCases{
		{
			tablePath: "color",
			expected: []string{
				"Blue", "Red", "Yellow",
			},
		},
		{
			tablePath: "places/country",
			expected: []string{
				"USA", "Mexico", "Canada",
			},
		},
		{
			tablePath: "places/castle/name",
			expected: []string{
				"Roogna", "Grayskull", "Castle AARrrrrgghhhh", "Edinburgh", "Neuschwanstein",
			},
		},
		{
			tablePath: "people/name",
			expected: []string{
				"Bob", "Sue",
			},
		},
		{
			tablePath: "things/item",
			expected: []string{
				"Dagger", "Coin", "Gem", "Sword",
			},
		},
		{
			tablePath: "things/fancy",
			expected: []string{
				"Shiny Dagger", "Shiny Coin", "Shiny Gem", "Shiny Sword",
			},
		},
		{
			tablePath: "nested/lookup",
			expected: []string{
				"Foo", "Bar", "Baz",
			},
		},
	}
	for i, tc := range tests {
		actual, err := tree.GetItem(tc.tablePath)
		if err != nil {
			t.Errorf("TestIndex: %v, Error: %v", i, err)
		}
		found := false
		for _, exepctedItem := range tc.expected {
			if actual == exepctedItem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find one of %s but got %s", tc.expected, actual)
		}
	}

}

func TestListTables(t *testing.T) {
	tree := NewTree()
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewRandomTableRenderer(tree), 1))),
	)
	if err := md.Convert(bytes.NewBufferString(testmakrkwon).Bytes(), &buf); err != nil {
		t.Error(err)
	}
	tests := []TestCases{
		{
			tablePath: "",
			expected: []string{
				"color", "nested/lookup", "nested/subnest/subtable", "nested/subnest/table", "nested/table", "places/country", "places/castle/name", "people/name", "things/item", "things/fancy",
			},
		},
	}
	for _, tc := range tests {
		actual := tree.ListTables(tc.tablePath)
		actualSorted := sort.StringSlice(actual)
		actualSorted.Sort()
		expectedSorted := sort.StringSlice(tc.expected)
		expectedSorted.Sort()
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Expected to find %s but got %s", expectedSorted, actualSorted)
		}
	}
}

const testmakrkwon = `| Color  |
| ------ |
| B----e |
| Red    |
| Yellow |

# Places

| Country |
| ------- |
| USA     |
| Mexico  |
| Canada  |

## Castle

| Name                 |
| -------------------- |
| Roogna               |
| Grayskull            |
| Castle AARrrrrgghhhh |
| Edinburgh            |
| Neuschwanstein       |

# People

| Name |
| ---- |
| Bob  |
| Sue  |

# Things

| Item                                              | 2d4 |
| ------------------------------------------------- | --- |
| Dagger                                            | 2   |
| Coin                                              | 3-6 |
| Gem                                               | 7   |
| Sword from Castle {{lookup "places/castle/name"}} | 8   |

| Fancy                          |
| ------------------------------ |
| Shiny {{lookup "things/item"}} |

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

`

package randomtable

import (
	"bytes"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const markdownfixture = "./tests/test.md"

type TestCases struct {
	tablePath string
	expected  []string
}

func TestHeaderLookups(t *testing.T) {
	tree := NewTree()
	source, err := ioutil.ReadFile(markdownfixture)
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewRandomTableRenderer(tree), 1))),
	)
	if err := md.Convert(source, &buf); err != nil {
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
	source, err := ioutil.ReadFile(markdownfixture)
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewRandomTableRenderer(tree), 1))),
	)
	if err := md.Convert(source, &buf); err != nil {
		t.Error(err)
	}
	tests := []TestCases{
		{
			tablePath: "",
			expected: []string{
				"color", "places/country", "places/castle/name", "people/name", "things/item", "things/fancy",
			},
		},
	}
	for _, tc := range tests {
		actual := tree.ListTables(tc.tablePath)
		actual = sort.StringSlice(actual)
		tc.expected = sort.StringSlice(tc.expected)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Expected to find %s but got %s", tc.expected, actual)
		}
	}
}

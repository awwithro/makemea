package cmd

import (
	"fmt"
	"makemea/randomtable"
	"sort"
)

func list(tree randomtable.Tree, prefix string) {
	sortedItems := sort.StringSlice(tree.ListTables(prefix))
	sortedItems.Sort()
	for _, item := range sortedItems {
		fmt.Println(item)
	}
}

package cmd

import (
	"fmt"
	"sort"

	"github.com/awwithro/makemea/randomtable"
)

func list(tree randomtable.Tree, prefix string) {
	sortedItems := sort.StringSlice(tree.ListTables(prefix))
	sortedItems.Sort()
	for _, item := range sortedItems {
		fmt.Println(item)
	}
}

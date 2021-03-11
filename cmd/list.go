package cmd

import (
	"fmt"

	"github.com/awwithro/makemea/randomtable"
)

func list(tree randomtable.Tree, prefix string) {

	for _, item := range tree.ListTables(prefix) {
		fmt.Println(item)
	}
}

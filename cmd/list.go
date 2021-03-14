package cmd

import (
	"fmt"

	"github.com/awwithro/makemea/randomtable"
	"github.com/spf13/cobra"
)

// ListAll is used to list hidden tables
var ListAll bool

func list(tree randomtable.Tree, prefix string, showHidden bool) {
	for _, item := range tree.ListTables(prefix, showHidden) {
		fmt.Println(item)
	}
}

var listCmd = &cobra.Command{
	Use:   "list [prefix]",
	Short: "list tables with the given prefix",
	Run: func(cmd *cobra.Command, args []string) {
		var tableName string
		if len(args) == 0 {
			tableName = ""
		} else {
			tableName = args[0]
		}
		tree := MustGetTree()
		tree.ValidateTables()
		list(tree, tableName, ListAll)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SetArgs([]string{""})
		}
		return nil
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&ListAll, "all", "a", false, "List hidden tables")
}

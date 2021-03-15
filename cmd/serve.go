package cmd

import (
	"github.com/awwithro/makemea/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [prefix]",
	Short: "runs a server that serves tables",
	Run: func(cmd *cobra.Command, args []string) {

		tree := MustGetTree()
		tree.ValidateTables()
		server := server.NewServer(tree)
		server.Run()
	},
}

// func init() {
// 	listCmd.PersistentFlags().BoolVarP(&ListAll, "all", "a", false, "List hidden tables")
// }

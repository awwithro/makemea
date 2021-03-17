package cmd

import (
	"github.com/awwithro/makemea/server"
	"github.com/spf13/cobra"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve [prefix]",
	Short: "runs a server that serves tables",
	Run: func(cmd *cobra.Command, args []string) {

		tree := MustGetTree()
		tree.ValidateTables()
		server := server.NewServer(tree)
		server.Run("0.0.0.0" + port)
	},
}

func init() {
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", ":8080", "Port for the server to listen on (:8181)")
}

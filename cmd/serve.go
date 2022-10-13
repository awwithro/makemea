package cmd

import (
	"log"
	"os"
	"time"

	"github.com/awwithro/makemea/server"
	"github.com/spf13/cobra"
	"golang.org/x/mod/sumdb/dirhash"
)

var port string
var addr string

var serveCmd = &cobra.Command{
	Use:   "serve [prefix]",
	Short: "runs a server that serves tables",
	Run:   serveCommand,
}

func init() {
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", ":8080", "Port for the server to listen on (:8181)")
	serveCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "127.0.0.1", "Address for the server to listen on (127.0.0.1)")
}

func serveCommand(cmd *cobra.Command, args []string) {
	tree := MustGetTree().WithHtmlFormatter()
	tree.ValidateTables()
	srv := server.NewServer(&tree)
	ticker := time.NewTicker(5 * time.Second)
	dir, _ := os.Getwd()
	hash, _ := dirhash.HashDir(dir, "", dirhash.DefaultHash)
	go func() {
		for {
			select {
			// Check for updated files every 5 seconds and reload the tree if things have changed
			case <-ticker.C:
				newHash, _ := dirhash.HashDir(dir, "", dirhash.DefaultHash)
				if hash != newHash {
					log.Print("Files have changed, reloading tables")
					newTree := MustGetTree()
					*&tree = newTree
					hash = newHash
				}
			}
		}
	}()
	srv.Run(addr + port)
}

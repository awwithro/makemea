package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"makemea/randomtable"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var rootCmd = &cobra.Command{
	Use:   "makemea",
	Short: "MakeMeA is a tool to let GMs roll on lookup tables composed in markdown",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		var tableName string
		if len(args) == 0 {
			tableName = ""
		} else {
			tableName = args[0]
		}
		tree := MustGetTree()

		ls, _ := cmd.Flags().GetBool("list")
		// list the tables
		if ls {
			fmt.Print(tree.ListTables(tableName))
		} else { // get an item from the table
			// table, err := tree.GetTable(tableName)
			// if err != nil {
			// 	log.Fatalf("No table exists with the name: %s", tableName)
			// }
			item, err := tree.GetItem(tableName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(item)
		}

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SetArgs([]string{""})
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseMarkdown(path string, tree randomtable.RandomTableTree) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(randomtable.NewRandomTableRenderer(tree), 1))),
	)
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
}

func loadTablesIntoTree(tree randomtable.RandomTableTree) error {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				parseMarkdown(path, tree)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func MustGetTree() randomtable.RandomTableTree {
	tree := randomtable.NewRandomTableTree()
	err := loadTablesIntoTree(tree)
	if err != nil {
		log.Fatalf("Unable to load tables: %v", err)
	}
	return tree
}

func init() {
	rootCmd.PersistentFlags().Bool("list", false, "list tables at the given prefix")

}

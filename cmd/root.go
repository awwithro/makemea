package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/awwithro/makemea/randomtable"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var rootCmd = &cobra.Command{
	Use:   "makemea <table_name>",
	Short: "MakeMeA is a tool to let GMs roll on lookup tables composed in markdown",
	Long: `MakeMeA is a tool to let GMs roll on lookup tables composed in markdown.
It will recursively search the current directory for any markdown files
and attempt to turn any tables in those files into tables that can be rolled on.`,
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		tree := MustGetTree()
		tree.ValidateTables()
		item, err := tree.GetItem(tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(item)

	},
	Args: cobra.MinimumNArgs(1),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseMarkdown(path string, tree randomtable.Tree) {
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

func loadTablesIntoTree(tree randomtable.Tree) error {
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

func MustGetTree() randomtable.Tree {
	tree := randomtable.NewTree()
	err := loadTablesIntoTree(tree)
	if err != nil {
		log.Fatalf("Unable to load tables: %v", err)
	}
	return tree
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(showCmd)
}

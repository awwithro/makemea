package cmd

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/awwithro/makemea/randomtable"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func show(tree randomtable.Tree, tableName string) {
	t, name, err := tree.GetTable(tableName)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(name, "/")
	shortName := s[len(s)-1]
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoFormatHeaders(false)
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(false)
	table.SetAutoWrapText(false)
	table.SetColWidth(1000)
	table.SetReflowDuringAutoWrap(false)
	table.SetRowLine(false)
	table = t.GetTable(table, strings.Title(shortName))
	table.Render()
}

var showCmd = &cobra.Command{
	Use:   "show [table]",
	Short: "Prints the contents of a table in markdown format",
	Run: func(cmd *cobra.Command, args []string) {
		var tableName string
		tableName = args[0]
		tree := MustGetTree()
		tree.ValidateTables()
		show(tree, tableName)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No Table specified to show")
		}
		return nil
	},
}

func init() {
}

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists the tables that have been loaded and parsed",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		tree := MustGetTree()
		tree.GetTable("test")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a table to roll on")
		}
		return nil
	},
}

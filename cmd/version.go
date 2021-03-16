package cmd

import "github.com/spf13/cobra"

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of makemea",
	Run: func(cmd *cobra.Command, args []string) {
		println("Version: ", Version)
	},
	Args: cobra.NoArgs,
}

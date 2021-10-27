package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var GoVersion = "none"
var CommitHash = "none"
var GitTag = "none"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of makemea",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoVersion: %v\nCommitHash: %v\nGitTag: %v\n", GoVersion, CommitHash, GitTag)
	},
	Args: cobra.NoArgs,
}

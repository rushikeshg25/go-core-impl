package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Show commit logs",
	Run: func(cmd *cobra.Command, args []string) {
		runCommit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func runCommit(cmd *cobra.Command, args []string) {
	fmt.Println("Showing commit logs")
}

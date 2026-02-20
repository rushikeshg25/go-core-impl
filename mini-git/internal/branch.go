package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to the staging area",
	Run: func(cmd *cobra.Command, args []string) {
		runBranch(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func runBranch(cmd *cobra.Command, args []string) {
	fmt.Println("Adding files to the staging area")
}

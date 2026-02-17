package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to the staging area",
	Run: func(cmd *cobra.Command, args []string) {
		runAdd(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	fmt.Println("Adding files to the staging area")
}

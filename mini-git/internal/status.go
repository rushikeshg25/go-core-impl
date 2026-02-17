package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Run: func(cmd *cobra.Command, args []string) {
		runStatus(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	fmt.Println("Showing the working tree status")
}

package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Run: func(cmd *cobra.Command, args []string) {
		runLog(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func runLog(cmd *cobra.Command, args []string) {
	fmt.Println("Showing commit logs")
}

package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout a commit",
	Run: func(cmd *cobra.Command, args []string) {
		runCheckout(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func runCheckout(cmd *cobra.Command, args []string) {
	fmt.Println("Checking out a commit")
}

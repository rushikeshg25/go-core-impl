package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repository",
	Run: func(cmd *cobra.Command, args []string) {
		runInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) {
	if _, err := os.Stat(".mini-git"); !os.IsNotExist(err) {
		fmt.Println("Mini-git repository already exists")
		return
	}

	dirs := []string{
		".mini-git/objects",
		".mini-git/refs/heads",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	headPath := filepath.Join(".mini-git", "HEAD")
	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		fmt.Printf("Error creating HEAD file: %v\n", err)
		return
	}

	fmt.Println("Initialized empty Mini-Git repository in .mini-git/")
}

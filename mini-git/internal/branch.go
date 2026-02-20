package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [name]",
	Short: "List or create branches",
	Run: func(_ *cobra.Command, args []string) {
		runBranch(args)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func runBranch(args []string) {
	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		fmt.Println("Not a mini-git repository (run 'mini-git init' first)")
		return
	}

	if len(args) == 0 {
		listBranches()
	} else {
		createBranch(args[0])
	}
}

func listBranches() {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		fmt.Printf("Error getting current branch: %v\n", err)
		return
	}

	files, err := os.ReadDir(".mini-git/refs/heads")
	if err != nil {
		fmt.Printf("Error reading branches: %v\n", err)
		return
	}

	for _, file := range files {
		prefix := "  "
		if file.Name() == currentBranch {
			prefix = "* "
		}
		fmt.Printf("%s%s\n", prefix, file.Name())
	}
}

func createBranch(name string) {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	branchPath := filepath.Join(".mini-git", "refs", "heads", currentBranch)
	commitHash, err := os.ReadFile(branchPath)
	if err != nil {
		fmt.Println("Error: Not a valid object name: 'main'. (Have you committed yet?)")
		return
	}

	newBranchPath := filepath.Join(".mini-git", "refs", "heads", name)
	if _, err := os.Stat(newBranchPath); err == nil {
		fmt.Printf("fatal: A branch named '%s' already exists.\n", name)
		return
	}

	if err := os.WriteFile(newBranchPath, commitHash, 0644); err != nil {
		fmt.Printf("Error creating branch: %v\n", err)
		return
	}
}

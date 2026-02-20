package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Run: func(_ *cobra.Command, _ []string) {
		runLog()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func runLog() {
	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		fmt.Println("Not a mini-git repository (run 'mini-git init' first)")
		return
	}

	commitHash, err := ResolveRef("HEAD")
	if err != nil {
		fmt.Println("fatal: your current branch 'main' does not have any commits yet")
		return
	}

	for commitHash != "" {
		data, err := ReadObject(commitHash)
		if err != nil {
			fmt.Printf("Error reading commit %s: %v\n", commitHash, err)
			break
		}

		lines := strings.Split(string(data), "\n")
		var author, message strings.Builder
		var parent string
		isMessage := false

		for _, line := range lines {
			if isMessage {
				message.WriteString(line)
				message.WriteString("\n")
				continue
			}
			if line == "" {
				isMessage = true
				continue
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}

			switch parts[0] {
			case "author":
				author.WriteString(parts[1])
			case "parent":
				parent = parts[1]
			}
		}

		fmt.Printf("\033[33mcommit %s\033[0m\n", commitHash)
		fmt.Printf("Author: %s\n", author.String())
		fmt.Println()
		fmt.Printf("    %s\n", strings.TrimSpace(message.String()))
		fmt.Println()

		commitHash = parent
	}
}

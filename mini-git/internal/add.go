package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [file...]",
	Short: "Add files to the staging area",
	Run: func(_ *cobra.Command, args []string) {
		runAdd(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(args []string) {
	if len(args) == 0 {
		fmt.Println("No files specified to add")
		return
	}

	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		fmt.Println("Not a mini-git repository (run 'mini-git init' first)")
		return
	}

	index := make(map[string]string)
	indexPath := ".mini-git/index"

	if _, err := os.Stat(indexPath); err == nil {
		file, err := os.Open(indexPath)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				parts := strings.Split(scanner.Text(), " ")
				if len(parts) == 2 {
					index[parts[1]] = parts[0]
				}
			}
			file.Close()
		}
	}

	for _, path := range args {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			continue
		}

		hash, err := StoreObject(data)
		if err != nil {
			fmt.Printf("Error storing object for %s: %v\n", path, err)
			continue
		}

		index[path] = hash
		fmt.Printf("Added %s\n", path)
	}

	file, err := os.Create(indexPath)
	if err != nil {
		fmt.Printf("Error creating index file: %v\n", err)
		return
	}
	defer file.Close()

	for path, hash := range index {
		fmt.Fprintf(file, "%s %s\n", hash, path)
	}
}

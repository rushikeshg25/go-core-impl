package internal

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Run: func(_ *cobra.Command, _ []string) {
		runStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus() {
	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		fmt.Println("Not a mini-git repository (run 'mini-git init' first)")
		return
	}

	index, err := loadIndex()
	if err != nil {
		fmt.Printf("Error loading index: %v\n", err)
		return
	}

	staged := []string{}
	modified := []string{}
	untracked := []string{}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".mini-git" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath := strings.TrimPrefix(path, "./")

		if hash, ok := index[relPath]; ok {
			currentHash, err := hashFile(path)
			if err != nil {
				return err
			}
			if currentHash != hash {
				modified = append(modified, relPath)
			} else {
				staged = append(staged, relPath)
			}
		} else {
			untracked = append(untracked, relPath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	currentBranch, err := getCurrentBranch()
	if err != nil {
		currentBranch = "unknown"
	}
	fmt.Printf("On branch %s\n", currentBranch)
	fmt.Println()

	if len(staged) > 0 {
		fmt.Println("Changes to be committed:")
		for _, file := range staged {
			fmt.Printf("\tnew file:   %s\n", file)
		}
		fmt.Println()
	}

	if len(modified) > 0 {
		fmt.Println("Changes not staged for commit:")
		for _, file := range modified {
			fmt.Printf("\tmodified:   %s\n", file)
		}
		fmt.Println()
	}

	if len(untracked) > 0 {
		fmt.Println("Untracked files:")
		for _, file := range untracked {
			fmt.Printf("\t%s\n", file)
		}
		fmt.Println()
	}

	if len(staged) == 0 && len(modified) == 0 && len(untracked) == 0 {
		fmt.Println("nothing to commit, working tree clean")
	}
}

func loadIndex() (map[string]string, error) {
	index := make(map[string]string)
	indexPath := ".mini-git/index"

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return index, nil
	}

	file, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 2 {
			index[parts[1]] = parts[0]
		}
	}

	return index, scanner.Err()
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

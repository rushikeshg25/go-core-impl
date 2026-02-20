package internal

import (
	"fmt"
	"os"
	"strings"
)

// getCurrentBranch returns the name of the current branch from HEAD
func getCurrentBranch() (string, error) {
	headPath := ".mini-git/HEAD"
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	content := string(data)
	if strings.HasPrefix(content, "ref: refs/heads/") {
		return strings.TrimSpace(strings.TrimPrefix(content, "ref: refs/heads/")), nil
	}

	return "", fmt.Errorf("HEAD is in detached state or invalid")
}

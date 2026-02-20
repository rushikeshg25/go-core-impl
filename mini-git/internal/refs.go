package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func ResolveRef(ref string) (string, error) {
	if ref == "HEAD" {
		branch, err := getCurrentBranch()
		if err != nil {
			return "", err
		}
		ref = branch
	}

	path := filepath.Join(".mini-git", "refs", "heads", ref)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

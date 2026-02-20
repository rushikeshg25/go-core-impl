package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit(t *testing.T) {
	tmpDir := t.TempDir()

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(oldCwd)
		if err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

	runInit(nil, nil)

	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		t.Errorf(".mini-git directory was not created")
	}

	subDirs := []string{".mini-git/objects", ".mini-git/refs/heads"}
	for _, dir := range subDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("directory %s was not created", dir)
		}
	}

	headPath := filepath.Join(".mini-git", "HEAD")
	content, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("failed to read HEAD file: %v", err)
	}
	expectedContent := "ref: refs/heads/main\n"
	if string(content) != expectedContent {
		t.Errorf("expected HEAD content %q, got %q", expectedContent, string(content))
	}
}

func TestRunInitAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runInit(nil, nil)
	runInit(nil, nil)

	if _, err := os.Stat(".mini-git"); os.IsNotExist(err) {
		t.Errorf(".mini-git directory disappeared after second init")
	}
}

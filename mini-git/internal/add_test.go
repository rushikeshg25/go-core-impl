package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAdd(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	// Initialize repo
	runInit(nil, nil)

	// Create a test file
	fileName := "test.txt"
	fileContent := "hello mini-git"
	if err := os.WriteFile(fileName, []byte(fileContent), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Run add command
	runAdd([]string{fileName})

	// Check if object exists
	// The hash of "hello mini-git" is 79347895493b827e8a939103c80a8c2041d50c7c
	// Let's verify it dynamically
	indexPath := ".mini-git/index"
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index: %v", err)
	}

	indexContent := string(data)
	if !strings.Contains(indexContent, fileName) {
		t.Errorf("index does not contain file name")
	}

	parts := strings.Split(indexContent, " ")
	if len(parts) < 1 {
		t.Fatalf("unexpected index format")
	}
	hash := parts[0]

	objectPath := filepath.Join(".mini-git", "objects", hash[:2], hash[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		t.Errorf("object file %s was not created", objectPath)
	}

	objectData, err := os.ReadFile(objectPath)
	if err != nil {
		t.Fatalf("failed to read object: %v", err)
	}
	if string(objectData) != fileContent {
		t.Errorf("expected object content %q, got %q", fileContent, string(objectData))
	}
}

func TestRunAddNoRepo(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	// Run add without init
	runAdd([]string{"test.txt"})

	if _, err := os.Stat(".mini-git/index"); !os.IsNotExist(err) {
		t.Errorf("index should not exist without repo")
	}
}

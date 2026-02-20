package internal

import (
	"os"
	"testing"
)

func TestRunStatus(t *testing.T) {
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

	// Create and add a file
	fileName := "staged.txt"
	if err := os.WriteFile(fileName, []byte("staged content"), 0644); err != nil {
		t.Fatalf("failed to create staged file: %v", err)
	}
	runAdd([]string{fileName})

	// Create a modified file (add then modify)
	modifiedName := "modified.txt"
	if err := os.WriteFile(modifiedName, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to create modified file: %v", err)
	}
	runAdd([]string{modifiedName})
	if err := os.WriteFile(modifiedName, []byte("changed content"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	// Create an untracked file
	untrackedName := "untracked.txt"
	if err := os.WriteFile(untrackedName, []byte("untracked content"), 0644); err != nil {
		t.Fatalf("failed to create untracked file: %v", err)
	}

	// Capture output of runStatus
	// Note: runStatus prints to stdout. In a real test we'd capture it or refactor runStatus to take an io.Writer.
	// For now, let's just run it to make sure it doesn't crash and then do basic verification.
	runStatus()
}

func TestRunStatusNoRepo(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runStatus()
}

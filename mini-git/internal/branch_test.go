package internal

import (
	"os"
	"testing"
)

func TestRunBranchList(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runInit(nil, nil)

	runBranch(nil)

	if err := os.WriteFile(".mini-git/refs/heads/main", []byte("somehash"), 0644); err != nil {
		t.Fatalf("failed to create main branch: %v", err)
	}

	runBranch(nil)
}

func TestRunBranchCreate(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runInit(nil, nil)

	runBranch([]string{"new-feature"})
	if _, err := os.Stat(".mini-git/refs/heads/new-feature"); err == nil {
		t.Errorf("branch should not have been created without a parent commit")
	}

	// Manually create a main branch file to simulate a commit
	if err := os.WriteFile(".mini-git/refs/heads/main", []byte("somehash"), 0644); err != nil {
		t.Fatalf("failed to create main branch: %v", err)
	}

	// Create branch
	runBranch([]string{"new-feature"})

	// Check if new branch exists
	data, err := os.ReadFile(".mini-git/refs/heads/new-feature")
	if err != nil {
		t.Fatalf("failed to read new branch: %v", err)
	}
	if string(data) != "somehash" {
		t.Errorf("expected branch to point to 'somehash', got %q", string(data))
	}
}

func TestRunBranchNoRepo(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runBranch(nil)
}

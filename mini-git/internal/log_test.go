package internal

import (
	"fmt"
	"os"
	"testing"
)

func TestRunLog(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runInit(nil, nil)

	c1Hash, _ := StoreObject([]byte("tree some-tree-hash\nauthor User <user@example.com>\n\nInitial commit\n"))

	c2Hash, _ := StoreObject([]byte(fmt.Sprintf("tree another-tree-hash\nparent %s\nauthor User <user@example.com>\n\nSecond commit\n", c1Hash)))

	if err := os.WriteFile(".mini-git/refs/heads/main", []byte(c2Hash), 0644); err != nil {
		t.Fatalf("failed to update main branch: %v", err)
	}
	runLog()
}

func TestRunLogNoCommits(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	runInit(nil, nil)
	runLog()
}

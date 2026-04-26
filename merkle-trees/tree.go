package main

import (
	"os"
	"path/filepath"
	"sort"
)

// BuildTree constructs a Merkle tree rooted at the given directory path.
// It recursively walks the filesystem, hashing file contents at leaves and
// computing directory hashes from sorted child hashes.
func BuildTree(rootPath string) (*Node, error) {
	// Resolve to absolute path for consistency
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "build", Path: absPath, Err: os.ErrInvalid}
	}

	return buildNode(absPath, info.Name(), "", nil)
}

// buildNode recursively creates a Merkle tree node for the given path.
//   - For files: reads content, hashes it, and returns a leaf node.
//   - For directories: recurses into children, sorts them, computes dir hash.
//
// Parameters:
//   - fullPath: absolute path to the file or directory
//   - name: basename of the entry
//   - relPath: relative path from the tree root (used for display/diff)
//   - parent: pointer to the parent node (nil for root)
func buildNode(fullPath, name, relPath string, parent *Node) (*Node, error) {
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name:   name,
		Path:   relPath,
		IsDir:  info.IsDir(),
		Parent: parent,
	}

	if !info.IsDir() {
		// Leaf node: hash file contents
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}
		node.Hash = hashData(data)
		return node, nil
	}

	// Internal node: read directory entries
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	// Sort entries by name for deterministic traversal
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		childName := entry.Name()

		// Skip ignored files/directories
		if ShouldIgnore(childName) {
			continue
		}

		childFull := filepath.Join(fullPath, childName)
		childRel := filepath.Join(relPath, childName)

		childNode, err := buildNode(childFull, childName, childRel, node)
		if err != nil {
			return nil, err
		}

		node.Children = append(node.Children, childNode)
	}

	// Compute directory hash from children
	node.Hash = hashDir(node)

	return node, nil
}

package main

import (
	"fmt"
	"sort"

	"github.com/zeebo/blake3"
)

// DefaultIgnorePatterns contains file/directory names that are skipped
// during tree construction (hidden files, VCS directories, OS artifacts).
var DefaultIgnorePatterns = map[string]bool{
	".git":       true,
	".gitignore": true,
	".DS_Store":  true,
	".hg":        true,
	".svn":       true,
	"node_modules": true,
}

// ShouldIgnore returns true if the given filename should be excluded from the tree.
func ShouldIgnore(name string) bool {
	return DefaultIgnorePatterns[name]
}

// hashData computes the BLAKE3-256 hash of raw byte content (used for files).
func hashData(data []byte) [32]byte {
	return blake3.Sum256(data)
}

// hashDir computes the BLAKE3-256 hash of a directory node by sorting its
// children by name and hashing the concatenation of "name:hex(hash)\n" pairs.
// This ensures the hash is deterministic regardless of filesystem enumeration order.
func hashDir(node *Node) [32]byte {
	// Sort children by name for deterministic hashing
	children := make([]*Node, len(node.Children))
	copy(children, node.Children)
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
	})

	// Build the payload: each child contributes "name:hex(hash)\n"
	payload := ""
	for _, child := range children {
		payload += fmt.Sprintf("%s:%x\n", child.Name, child.Hash)
	}

	return blake3.Sum256([]byte(payload))
}

// sortNodesByName sorts a slice of Nodes in-place by their Name field.
func sortNodesByName(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
}

package main

import (
	"fmt"
	"strings"
)

// Difference describes a single change detected between two Merkle trees.
type Difference struct {
	Path    string // relative path of the changed node
	Type    string // "modified", "added", or "deleted"
	OldHash string // hex-encoded hash from tree A (empty for "added")
	NewHash string // hex-encoded hash from tree B (empty for "deleted")
}

// CompareTrees walks two Merkle trees in parallel and returns a list of
// differences. It reports files/directories that were added, deleted, or
// modified between tree A and tree B.
func CompareTrees(a, b *Node) []Difference {
	var diffs []Difference
	compareTrees(a, b, &diffs)
	return diffs
}

// compareTrees is the recursive helper for CompareTrees.
func compareTrees(a, b *Node, diffs *[]Difference) {
	// If both hashes match, the entire subtree is identical — skip.
	if a.Hash == b.Hash {
		return
	}

	// Build child maps for O(1) lookup
	aChildren := childMap(a)
	bChildren := childMap(b)

	// Check for deleted or modified entries (in A but possibly changed/missing in B)
	for name, aChild := range aChildren {
		bChild, exists := bChildren[name]
		if !exists {
			// Present in A, missing in B → deleted
			*diffs = append(*diffs, Difference{
				Path:    aChild.Path,
				Type:    "deleted",
				OldHash: fmt.Sprintf("%x", aChild.Hash),
			})
			continue
		}

		if aChild.Hash != bChild.Hash {
			if aChild.IsDir && bChild.IsDir {
				// Both are directories — recurse deeper
				compareTrees(aChild, bChild, diffs)
			} else {
				// File was modified (or type changed, e.g., file→dir)
				*diffs = append(*diffs, Difference{
					Path:    aChild.Path,
					Type:    "modified",
					OldHash: fmt.Sprintf("%x", aChild.Hash),
					NewHash: fmt.Sprintf("%x", bChild.Hash),
				})
			}
		}
	}

	// Check for added entries (in B but not in A)
	for name, bChild := range bChildren {
		if _, exists := aChildren[name]; !exists {
			*diffs = append(*diffs, Difference{
				Path:    bChild.Path,
				Type:    "added",
				NewHash: fmt.Sprintf("%x", bChild.Hash),
			})
		}
	}
}

// childMap builds a lookup map from child name to *Node.
func childMap(n *Node) map[string]*Node {
	m := make(map[string]*Node, len(n.Children))
	for _, child := range n.Children {
		m[child.Name] = child
	}
	return m
}

// PrintTree pretty-prints the Merkle tree to stdout with indentation.
// Each node shows its name and a truncated hash.
func PrintTree(node *Node, indent int) {
	prefix := strings.Repeat("  ", indent)
	hashHex := fmt.Sprintf("%x", node.Hash)

	// Truncate hash to first 16 hex chars for readability
	shortHash := hashHex
	if len(shortHash) > 16 {
		shortHash = shortHash[:16]
	}

	if node.IsDir {
		fmt.Printf("%s📁 %s/  [%s…]\n", prefix, node.Name, shortHash)
		for _, child := range node.Children {
			PrintTree(child, indent+1)
		}
	} else {
		fmt.Printf("%s📄 %s   [%s…]\n", prefix, node.Name, shortHash)
	}
}

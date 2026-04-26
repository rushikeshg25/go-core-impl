package main

import (
	"fmt"
	"os"
)

// Node represents a single node in the Merkle tree.
// Leaf nodes correspond to files; internal nodes correspond to directories.
type Node struct {
	Name     string   // file or directory basename
	Path     string   // full relative path from the tree root
	IsDir    bool     // true if this node represents a directory
	Hash     [32]byte // BLAKE3 hash of contents (file) or children (dir)
	Children []*Node  // child nodes; non-nil only for directories
	Parent   *Node    // back-pointer to parent; nil for root
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "hash":
		cmdHash()
	case "print":
		cmdPrint()
	case "diff":
		cmdDiff()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: merkle-tree <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  hash  <dir>          Build Merkle tree and print root hash")
	fmt.Fprintln(os.Stderr, "  print <dir>          Build Merkle tree and pretty-print the tree")
	fmt.Fprintln(os.Stderr, "  diff  <dir1> <dir2>  Compare two directory trees and show differences")
}

// cmdHash builds a Merkle tree for the given directory and prints the root hash.
func cmdHash() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: merkle-tree hash <dir>")
		os.Exit(1)
	}
	dir := os.Args[2]

	root, err := BuildTree(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%x\n", root.Hash)
}

// cmdPrint builds a Merkle tree and pretty-prints it.
func cmdPrint() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: merkle-tree print <dir>")
		os.Exit(1)
	}
	dir := os.Args[2]

	root, err := BuildTree(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	PrintTree(root, 0)
}

// cmdDiff builds two Merkle trees and reports differences.
func cmdDiff() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: merkle-tree diff <dir1> <dir2>")
		os.Exit(1)
	}
	dir1 := os.Args[2]
	dir2 := os.Args[3]

	treeA, err := BuildTree(dir1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building tree for %s: %v\n", dir1, err)
		os.Exit(1)
	}

	treeB, err := BuildTree(dir2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building tree for %s: %v\n", dir2, err)
		os.Exit(1)
	}

	diffs := CompareTrees(treeA, treeB)
	if len(diffs) == 0 {
		fmt.Println("Trees are identical.")
		return
	}

	fmt.Printf("Found %d difference(s):\n\n", len(diffs))
	for _, d := range diffs {
		switch d.Type {
		case "added":
			fmt.Printf("  + %-10s %s\n", d.Type, d.Path)
		case "deleted":
			fmt.Printf("  - %-10s %s\n", d.Type, d.Path)
		case "modified":
			fmt.Printf("  ~ %-10s %s\n", d.Type, d.Path)
			fmt.Printf("    old: %s\n", d.OldHash)
			fmt.Printf("    new: %s\n", d.NewHash)
		}
	}
}

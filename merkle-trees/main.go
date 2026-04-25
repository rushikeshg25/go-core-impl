package main

import "fmt"

type Node struct {
	child    []*Node
	fileName string
	path     string
	parent   *Node
	isDir    bool
	hash     string
}

func main() {
	fmt.Println("Hello World")
}

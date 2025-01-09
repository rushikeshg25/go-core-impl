package main

import "fmt"

func main() {
	ch := NewConsistentHash(50)
	node1 := &StorageNode{Host: "node1"}
	node2 := &StorageNode{Host: "node2"}

	ch.AddNode(node1)
	ch.AddNode(node2)

	item := "item1"
	assignedNode, _ := ch.Assign(item)
	fmt.Printf("Item %s is assigned to node %s\n", item, assignedNode.Host)

	ch.Plot(item, node1)
}

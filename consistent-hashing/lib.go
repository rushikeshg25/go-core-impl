package main

type Chash struct {
	Nodes []Node
}

type Node struct {
	Host string
}

func (c *Chash) AddNode(host string) {}

func (c *Chash) RemoveNode(host string) {}

func (c *Chash) GetNode(key string)  {}

func (c *Chash) Hash_fn(){}
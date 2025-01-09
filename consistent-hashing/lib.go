package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

type ConsistentHash struct {
	Nodes []*StorageNode
	Keys  []int
	totalSlots int
}

type StorageNode struct {
	Host string
}

func NewConsistentHash(totalSlots int) *ConsistentHash {
	return &ConsistentHash{
		Keys:       []int{},
		Nodes:      []*StorageNode{},
		totalSlots: totalSlots,
	}
}

func (c *ConsistentHash) AddNode(node *StorageNode) (int,error){
	if(len(c.Nodes)==c.totalSlots){
		return 0,errors.New("No more slots available")
	}
	key:=c.Hash_fn(node.Host,c.totalSlots)
	index:=sort.SearchInts(c.Keys,key)

	if index < len(c.Keys) && c.Keys[index] == key {
		return 0, errors.New("collision occurred")
	}

	c.Keys = append(c.Keys[:index], append([]int{key}, c.Keys[index:]...)...)
	c.Nodes = append(c.Nodes[:index], append([]*StorageNode{node}, c.Nodes[index:]...)...)

	return key, nil
}

func (ch *ConsistentHash) RemoveNode(node *StorageNode) (int, error) {
	if len(ch.Keys) == 0 {
		return 0, errors.New("hash space is empty")
	}

	key := ch.Hash_fn(node.Host, ch.totalSlots)
	index := sort.SearchInts(ch.Keys, key)

	
	if index >= len(ch.Keys) || ch.Keys[index] != key {
		return 0, errors.New("node does not exist")
	}

	
	ch.Keys = append(ch.Keys[:index], ch.Keys[index+1:]...)
	ch.Nodes = append(ch.Nodes[:index], ch.Nodes[index+1:]...)

	return key, nil
}

func (ch *ConsistentHash) Assign(item string) (*StorageNode, error) {
	if len(ch.Keys) == 0 {
		return nil, errors.New("hash ring is empty")
	}

	key := ch.Hash_fn(item, ch.totalSlots)
	index := sort.SearchInts(ch.Keys, key)

	
	if index == len(ch.Keys) {
		index = 0
	}

	return ch.Nodes[index], nil
}

func (ch *ConsistentHash) Plot(item string, node *StorageNode) {
	fmt.Println("Hash Ring Visualization:")
	fmt.Println("Keys:", ch.Keys)
	if item != "" {
		fmt.Printf("Item: %s, Hash: %d\n", item, ch.Hash_fn(item, ch.totalSlots))
	}
	if node != nil {
		fmt.Printf("Node: %s, Hash: %d\n", node.Host, ch.Hash_fn(node.Host, ch.totalSlots))
	}
}

func (c *ConsistentHash) Hash_fn(key string,totalSlots int)int{
	hash := sha256.Sum256([]byte(key))
	return int(binary.BigEndian.Uint64(hash[:8])) % totalSlots
}
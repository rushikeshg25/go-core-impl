package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rushikeshg25/raft/internal/raft"
	"github.com/rushikeshg25/raft/internal/rpc"
)

func main() {
	id := flag.Int("id", 0, "Node ID")
	cluster := flag.String("cluster", "localhost:8000,localhost:8001,localhost:8002", "Comma-separated cluster addresses")
	flag.Parse()

	peers := []string{}
	current := ""
	for _, char := range *cluster {
		if char == ',' {
			peers = append(peers, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	peers = append(peers, current)

	if *id < 0 || *id >= len(peers) {
		log.Fatalf("Invalid node ID %d for cluster of size %d", *id, len(peers))
	}

	address := peers[*id]

	node := raft.NewRaft(*id, peers, rpc.Call)
	server := rpc.NewServer(node)

	log.Printf("Starting Node %d on %s", *id, address)
	if err := server.Start(address); err != nil {
		log.Fatalf("Failed to start RPC server: %v", err)
	}

	node.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Printf("Shutting down Node %d", *id)
}

package main

import (
	"flag"
	"fmt"
	"github.com/rushikeshg25/raft/internal/raft"
	"github.com/rushikeshg25/raft/internal/rpc"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	id := flag.Int("id", 0, "node ID")
	cluster := flag.String("cluster", "localhost:8000,localhost:8001,localhost:8002", "static ordered peer addresses")
	state := flag.String("state", "", "persistent state file")
	submit := flag.String("submit", "", "send a command to this leader address")
	command := flag.String("command", "", "command bytes")
	status := flag.String("status", "", "read status from this address")
	flag.Parse()
	if *status != "" {
		var reply raft.StatusReply
		if !rpc.Call(*status, "Raft.Status", &raft.StatusArgs{}, &reply) {
			log.Fatal("status unavailable")
		}
		fmt.Printf("role=%d term=%d commit=%d last=%d\n", reply.Role, reply.Term, reply.CommitIndex, reply.LastLogIndex)
		return
	}
	if *submit != "" {
		var reply raft.SubmitReply
		if !rpc.Call(*submit, "Raft.Submit", &raft.SubmitArgs{Command: []byte(*command)}, &reply) {
			log.Fatal("proposal rejected or RPC timed out; outcome may be unknown")
		}
		fmt.Printf("appended index=%d term=%d; check commit status before treating it as committed\n", reply.Index, reply.Term)
		return
	}
	if *state == "" {
		*state = fmt.Sprintf("data/node-%d/state.json", *id)
	}
	peers := strings.Split(*cluster, ",")
	node, e := raft.NewPersistentRaft(*id, peers, rpc.Call, *state)
	if e != nil {
		log.Fatal(e)
	}
	server := rpc.NewServer(node)
	if e = server.Start(peers[*id]); e != nil {
		log.Fatal(e)
	}
	node.Start()
	log.Printf("node %d listening on %s", *id, peers[*id])
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	server.Stop()
	node.Stop()
}

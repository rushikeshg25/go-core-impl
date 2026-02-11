package rpc

import (
	"log"
	"net"
	"net/rpc"

	"github.com/rushikeshg25/raft/internal/raft"
)

type Server struct {
	node *raft.Raft
}

func NewServer(node *raft.Raft) *Server {
	return &Server{node: node}
}

func (s *Server) Start(address string) error {
	rpcServer := rpc.NewServer()
	err := rpcServer.RegisterName("Raft", s.node)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("[Server] Listening on %s", address)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[Server] Accept error: %v", err)
				continue
			}
			go rpcServer.ServeConn(conn)
		}
	}()

	return nil
}

func Call(address string, method string, args interface{}, reply interface{}) bool {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return false
	}
	defer client.Close()

	err = client.Call(method, args, reply)
	if err != nil {
		log.Printf("[RPC] Call error %s: %v", method, err)
		return false
	}

	return true
}

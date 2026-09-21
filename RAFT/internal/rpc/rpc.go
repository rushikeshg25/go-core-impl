package rpc

import (
	"errors"
	"github.com/rushikeshg25/raft/internal/raft"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Server struct {
	node     *raft.Raft
	mu       sync.Mutex
	listener net.Listener
	conns    map[net.Conn]bool
	stopped  bool
	wg       sync.WaitGroup
}

func NewServer(node *raft.Raft) *Server { return &Server{node: node, conns: map[net.Conn]bool{}} }
func (s *Server) Start(address string) error {
	l, e := net.Listen("tcp", address)
	if e != nil {
		return e
	}
	if e = s.Serve(l); e != nil {
		l.Close()
	}
	return e
}
func (s *Server) Serve(l net.Listener) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped || s.listener != nil {
		return errors.New("server already started or stopped")
	}
	server := rpc.NewServer()
	if e := server.RegisterName("Raft", s.node); e != nil {
		return e
	}
	s.listener = l
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			c, e := l.Accept()
			if e != nil {
				return
			}
			s.mu.Lock()
			if s.stopped {
				s.mu.Unlock()
				c.Close()
				return
			}
			s.conns[c] = true
			s.wg.Add(1)
			s.mu.Unlock()
			go func() {
				defer s.wg.Done()
				defer func() { c.Close(); s.mu.Lock(); delete(s.conns, c); s.mu.Unlock() }()
				server.ServeConn(c)
			}()
		}
	}()
	return nil
}
func (s *Server) Stop() {
	s.mu.Lock()
	s.stopped = true
	if s.listener != nil {
		s.listener.Close()
	}
	for c := range s.conns {
		c.Close()
	}
	s.mu.Unlock()
	s.wg.Wait()
}
func Call(address, method string, args, reply interface{}) bool {
	conn, e := net.DialTimeout("tcp", address, 200*time.Millisecond)
	if e != nil {
		return false
	}
	conn.SetDeadline(time.Now().Add(300 * time.Millisecond))
	client := rpc.NewClient(conn)
	defer client.Close()
	return client.Call(method, args, reply) == nil
}

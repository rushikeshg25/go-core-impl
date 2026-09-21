package rpc

import (
	"github.com/rushikeshg25/raft/internal/raft"
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestRealClusterElectionAndReplication(t *testing.T) {
	peers := make([]string, 3)
	listeners := make([]net.Listener, 3)
	for i := range peers {
		l, e := net.Listen("tcp", "127.0.0.1:0")
		if e != nil {
			t.Fatal(e)
		}
		listeners[i] = l
		peers[i] = l.Addr().String()
		defer l.Close()
	}
	nodes := make([]*raft.Raft, 3)
	for i := range nodes {
		r, e := raft.NewPersistentRaft(i, peers, Call, filepath.Join(t.TempDir(), "state"))
		if e != nil {
			t.Fatal(e)
		}
		nodes[i] = r
		s := NewServer(r)
		if e = s.Serve(listeners[i]); e != nil {
			t.Fatal(e)
		}
		r.Start()
		defer r.Stop()
		defer s.Stop()
	}
	deadline := time.Now().Add(5 * time.Second)
	leader := -1
	for time.Now().Before(deadline) && leader < 0 {
		for i, p := range peers {
			var status raft.StatusReply
			if Call(p, "Raft.Status", &raft.StatusArgs{}, &status) && status.Role == raft.Leader {
				leader = i
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	if leader < 0 {
		t.Fatal("no leader")
	}
	var reply raft.SubmitReply
	if !Call(peers[leader], "Raft.Submit", &raft.SubmitArgs{Command: []byte("hello")}, &reply) {
		t.Fatal("proposal rejected")
	}
	for time.Now().Before(deadline) {
		all := true
		for _, r := range nodes {
			got := r.Applied()
			all = all && len(got) == 1
			if len(got) > 0 && string(got[0].Command) != "hello" {
				t.Fatal("wrong command")
			}
		}
		if all {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("not replicated to every node")
}

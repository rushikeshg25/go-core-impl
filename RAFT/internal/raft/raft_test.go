package raft

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type testNetwork struct {
	mu      sync.Mutex
	nodes   []*Raft
	blocked map[int]bool
}

func (n *testNetwork) sender(from int) func(string, string, interface{}, interface{}) bool {
	return func(addr, method string, args, reply interface{}) bool {
		n.mu.Lock()
		id := int(addr[0] - '0')
		r := n.nodes[id]
		blocked := n.blocked[from] || n.blocked[id]
		n.mu.Unlock()
		if blocked {
			return false
		}
		switch method {
		case "Raft.RequestVote":
			return r.RequestVote(args.(*RequestVoteArgs), reply.(*RequestVoteReply)) == nil
		case "Raft.AppendEntries":
			return r.AppendEntries(args.(*AppendEntriesArgs), reply.(*AppendEntriesReply)) == nil
		}
		return false
	}
}
func wait(t *testing.T, f func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if f() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition timed out")
}
func role(r *Raft) NodeRole { r.mu.Lock(); defer r.mu.Unlock(); return r.role }
func TestPartitionRepairAndRestart(t *testing.T) {
	n := &testNetwork{nodes: make([]*Raft, 3), blocked: map[int]bool{}}
	paths := make([]string, 3)
	peers := []string{"0", "1", "2"}
	for i := range n.nodes {
		paths[i] = filepath.Join(t.TempDir(), "state.json")
		r, e := NewPersistentRaft(i, peers, n.sender(i), paths[i])
		if e != nil {
			t.Fatal(e)
		}
		n.nodes[i] = r
		defer r.Stop()
	}
	leader := n.nodes[0]
	leader.startElection()
	wait(t, func() bool { return role(leader) == Leader })
	input := []byte("first")
	if _, e := leader.Propose(input); e != nil {
		t.Fatal(e)
	}
	input[0] = 'X'
	wait(t, func() bool { leader.replicate(); return len(n.nodes[2].Applied()) == 1 })
	if string(n.nodes[2].Applied()[0].Command) != "first" {
		t.Fatal("command alias")
	}
	n.mu.Lock()
	n.blocked[0] = true
	n.mu.Unlock()
	orphan, e := leader.Propose([]byte("orphan"))
	if e != nil {
		t.Fatal(e)
	}
	leader.replicate()
	time.Sleep(20 * time.Millisecond)
	leader.mu.Lock()
	committed := leader.commitIndex
	leader.mu.Unlock()
	if committed >= orphan {
		t.Fatal("minority committed")
	}
	next := n.nodes[1]
	next.startElection()
	wait(t, func() bool { return role(next) == Leader })
	next.Propose([]byte("second"))
	wait(t, func() bool { next.replicate(); return len(n.nodes[2].Applied()) == 2 })
	n.mu.Lock()
	n.blocked[0] = false
	n.mu.Unlock()
	wait(t, func() bool { next.replicate(); return len(leader.Applied()) == 2 })
	for _, r := range n.nodes {
		got := r.Applied()
		if len(got) != 2 || string(got[0].Command) != "first" || string(got[1].Command) != "second" {
			t.Fatal(got)
		}
		r.Stop()
	}
	restarted, e := NewPersistentRaft(0, peers, n.sender(0), paths[0])
	if e != nil {
		t.Fatal(e)
	}
	defer restarted.Stop()
	if len(restarted.Applied()) != 2 {
		t.Fatal("lost committed prefix")
	}
}
func TestVoteDurabilityAndFreshness(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	send := func(string, string, interface{}, interface{}) bool { return false }
	r, e := NewPersistentRaft(0, []string{"0", "1", "2"}, send, path)
	if e != nil {
		t.Fatal(e)
	}
	var reply RequestVoteReply
	r.RequestVote(&RequestVoteArgs{Term: 2, CandidateId: 1}, &reply)
	if !reply.VoteGranted {
		t.Fatal("vote rejected")
	}
	r.Stop()
	r, e = NewPersistentRaft(0, []string{"0", "1", "2"}, send, path)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Stop()
	r.RequestVote(&RequestVoteArgs{Term: 2, CandidateId: 2}, &reply)
	if reply.VoteGranted {
		t.Fatal("double vote after restart")
	}
	var a AppendEntriesReply
	if e = r.AppendEntries(&AppendEntriesArgs{Term: 3, LeaderId: 1, Entries: []LogEntry{{Term: 3, Command: []byte("x")}}, LeaderCommit: 1}, &a); e != nil || !a.Success {
		t.Fatal(a, e)
	}
	r.RequestVote(&RequestVoteArgs{Term: 4, CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0}, &reply)
	if reply.VoteGranted {
		t.Fatal("stale candidate elected")
	}
	if e = r.AppendEntries(&AppendEntriesArgs{Term: 5, LeaderId: 1, Entries: []LogEntry{{Term: 5}}}, &a); e == nil {
		t.Fatal("committed prefix replaced")
	}
}
func TestSingleNodeAndPersistenceFailure(t *testing.T) {
	send := func(string, string, interface{}, interface{}) bool { return false }
	r, _ := NewPersistentRaft(0, []string{"0"}, send, filepath.Join(t.TempDir(), "state"))
	defer r.Stop()
	r.startElection()
	if _, e := r.Propose([]byte("one")); e != nil || len(r.Applied()) != 1 {
		t.Fatal(e, r.Applied())
	}
	r.path = t.TempDir()
	if _, e := r.Propose([]byte("two")); e == nil {
		t.Fatal("acknowledged failed persistence")
	}
	var vote RequestVoteReply
	if r.RequestVote(&RequestVoteArgs{Term: 10, CandidateId: 0}, &vote) == nil || vote.VoteGranted {
		t.Fatal("poisoned node voted")
	}
}

func TestOrderedApplicationCheckpoint(t *testing.T) {
	r := NewRaft(0, []string{"0"}, func(string, string, interface{}, interface{}) bool { return false })
	defer r.Stop()
	r.startElection()
	r.Propose([]byte("a"))
	r.Propose([]byte("b"))
	count := 0
	failure := errors.New("application failed")
	checkpoint, e := r.ApplyTo(0, func(entry AppliedEntry) error {
		count++
		var status StatusReply
		if err := r.Status(&StatusArgs{}, &status); err != nil {
			t.Fatal(err)
		}
		if string(entry.Command) == "b" {
			return failure
		}
		return nil
	})
	if !errors.Is(e, failure) || count != 2 || checkpoint != 2 {
		t.Fatal(checkpoint, count, e)
	}
	checkpoint, e = r.ApplyTo(checkpoint, func(entry AppliedEntry) error {
		if string(entry.Command) != "b" {
			t.Fatal("reapplied acknowledged entry")
		}
		return nil
	})
	if e != nil || checkpoint != 3 {
		t.Fatal(checkpoint, e)
	}
}

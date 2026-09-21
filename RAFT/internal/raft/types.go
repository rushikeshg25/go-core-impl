package raft

import (
	"sync"
	"time"
)

type NodeRole int

const (
	Follower NodeRole = iota
	Candidate
	Leader
)

type LogEntry struct {
	Term    int
	Command []byte
	Noop    bool
}
type Raft struct {
	mu                      sync.Mutex
	peers                   []string
	id                      int
	currentTerm, votedFor   int
	log                     []LogEntry
	commitIndex             int
	nextIndex, matchIndex   []int
	inflight                []bool
	role                    NodeRole
	heartbeat, election     time.Duration
	deadline, lastHeartbeat time.Time
	sendRPC                 func(string, string, interface{}, interface{}) bool
	path                    string
	err                     error
	started, stopped        bool
	stop                    chan struct{}
	wg                      sync.WaitGroup
}
type RequestVoteArgs struct{ Term, CandidateId, LastLogIndex, LastLogTerm int }
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}
type AppendEntriesArgs struct {
	Term, LeaderId, PrevLogIndex, PrevLogTerm int
	Entries                                   []LogEntry
	LeaderCommit                              int
}
type AppendEntriesReply struct {
	Term      int
	Success   bool
	NextIndex int
}
type SubmitArgs struct{ Command []byte }
type SubmitReply struct {
	Index, Term int
	Leader      bool
}
type StatusArgs struct{}
type StatusReply struct {
	Term, CommitIndex, LastLogIndex int
	Role                            NodeRole
}
type AppliedEntry struct {
	Index   int
	Command []byte
}

func NewRaft(id int, peers []string, send func(string, string, interface{}, interface{}) bool) *Raft {
	r, e := NewPersistentRaft(id, peers, send, "")
	if e != nil {
		panic(e)
	}
	return r
}

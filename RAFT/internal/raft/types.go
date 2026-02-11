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
	Command interface{}
}

type Raft struct {
	mu    sync.Mutex
	peers []string
	id    int

	currentTerm int
	votedFor    int
	log         []LogEntry

	commitIndex int
	lastApplied int

	nextIndex  []int
	matchIndex []int

	role        NodeRole
	heartbeat   time.Duration
	election    time.Duration
	lastContact time.Time
	sendRPC     func(address string, method string, args interface{}, reply interface{}) bool
}

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func NewRaft(id int, peers []string, sendRPC func(string, string, interface{}, interface{}) bool) *Raft {
	return &Raft{
		id:          id,
		peers:       peers,
		role:        Follower,
		votedFor:    -1,
		currentTerm: 0,
		heartbeat:   100 * time.Millisecond,
		election:    300 * time.Millisecond,
		log:         make([]LogEntry, 0),
		sendRPC:     sendRPC,
	}
}

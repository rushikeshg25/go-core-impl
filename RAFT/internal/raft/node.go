package raft

import (
	"log"
	"math/rand"
	"time"
)

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	log.Printf("[Node %d] Received RequestVote from %d (Term: %d)", rf.id, args.CandidateId, args.Term)

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return nil
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.role = Follower
		rf.votedFor = -1
	}

	reply.Term = rf.currentTerm

	canVote := rf.votedFor == -1 || rf.votedFor == args.CandidateId

	if canVote {
		rf.votedFor = args.CandidateId
		rf.lastContact = time.Now()
		reply.VoteGranted = true
		log.Printf("[Node %d] Voted for %d in Term %d", rf.id, args.CandidateId, rf.currentTerm)
	} else {
		reply.VoteGranted = false
	}

	return nil
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return nil
	}

	rf.lastContact = time.Now()

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.role = Follower
		rf.votedFor = -1
	} else if rf.role == Candidate {
		rf.role = Follower
	}

	reply.Term = rf.currentTerm
	reply.Success = true

	return nil
}

func (rf *Raft) Start() {
	rf.mu.Lock()
	rf.lastContact = time.Now()
	rf.mu.Unlock()
	go rf.ticker()
}

func (rf *Raft) ticker() {
	for {
		rf.mu.Lock()
		role := rf.role
		lastContact := rf.lastContact
		rf.mu.Unlock()

		switch role {
		case Follower, Candidate:
			timeout := rf.randomElectionTimeout()
			if time.Since(lastContact) > timeout {
				rf.startElection()
			}
			time.Sleep(20 * time.Millisecond)
		case Leader:
			rf.sendHeartbeats()
			time.Sleep(rf.heartbeat)
		}
	}
}

func (rf *Raft) randomElectionTimeout() time.Duration {
	r := rand.Intn(150)
	return rf.election + time.Duration(r)*time.Millisecond
}

func (rf *Raft) startElection() {
	rf.mu.Lock()
	rf.role = Candidate
	rf.currentTerm++
	rf.votedFor = rf.id
	rf.lastContact = time.Now()
	term := rf.currentTerm
	peers := rf.peers
	id := rf.id
	rf.mu.Unlock()

	log.Printf("[Node %d] Starting election for Term %d", rf.id, term)

	votes := 1
	for i, addr := range peers {
		if i == id {
			continue
		}

		go func(peerAddr string) {
			args := RequestVoteArgs{
				Term:        term,
				CandidateId: id,
			}
			var reply RequestVoteReply
			if rf.sendRPC(peerAddr, "Raft.RequestVote", &args, &reply) {
				rf.mu.Lock()
				defer rf.mu.Unlock()

				if rf.role != Candidate || rf.currentTerm != term {
					return
				}

				if reply.Term > rf.currentTerm {
					rf.currentTerm = reply.Term
					rf.role = Follower
					rf.votedFor = -1
					return
				}

				if reply.VoteGranted {
					votes++
					if votes > len(peers)/2 {
						log.Printf("[Node %d] Became LEADER for Term %d", rf.id, rf.currentTerm)
						rf.role = Leader
						rf.nextIndex = make([]int, len(peers))
						rf.matchIndex = make([]int, len(peers))
						for i := range rf.nextIndex {
							rf.nextIndex[i] = len(rf.log)
						}
					}
				}
			}
		}(addr)
	}
}

func (rf *Raft) sendHeartbeats() {
	rf.mu.Lock()
	term := rf.currentTerm
	peers := rf.peers
	id := rf.id
	rf.mu.Unlock()

	for i, addr := range peers {
		if i == id {
			continue
		}

		go func(peerAddr string) {
			args := AppendEntriesArgs{
				Term:     term,
				LeaderId: id,
			}
			var reply AppendEntriesReply
			if rf.sendRPC(peerAddr, "Raft.AppendEntries", &args, &reply) {
				rf.mu.Lock()
				defer rf.mu.Unlock()

				if rf.role != Leader || rf.currentTerm != term {
					return
				}

				if reply.Term > rf.currentTerm {
					rf.currentTerm = reply.Term
					rf.role = Follower
					rf.votedFor = -1
				}
			}
		}(addr)
	}
}

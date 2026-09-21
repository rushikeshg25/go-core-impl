package raft

import (
	"bytes"
	"errors"
	"math/rand"
	"time"
)

func (r *Raft) resetDeadline() {
	r.deadline = time.Now().Add(r.election + time.Duration(rand.Int63n(int64(r.election/2))))
}
func (r *Raft) stepDown(term int) error {
	r.role = Follower
	if term > r.currentTerm {
		r.currentTerm = term
		r.votedFor = -1
		return r.persist()
	}
	return nil
}
func (r *Raft) RequestVote(a *RequestVoteArgs, b *RequestVoteReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e := r.check(); e != nil {
		return e
	}
	if a.CandidateId < 0 || a.CandidateId >= len(r.peers) || a.Term < 0 || a.LastLogIndex < 0 || a.LastLogTerm < 0 {
		return errors.New("invalid vote request")
	}
	if a.Term > r.currentTerm {
		if e := r.stepDown(a.Term); e != nil {
			return e
		}
	}
	b.Term = r.currentTerm
	b.VoteGranted = false
	if a.Term < r.currentTerm {
		return nil
	}
	last := len(r.log) - 1
	upToDate := a.LastLogTerm > r.log[last].Term || a.LastLogTerm == r.log[last].Term && a.LastLogIndex >= last
	if upToDate && (r.votedFor == -1 || r.votedFor == a.CandidateId) {
		r.votedFor = a.CandidateId
		if e := r.persist(); e != nil {
			return e
		}
		r.resetDeadline()
		b.VoteGranted = true
	}
	return nil
}
func (r *Raft) AppendEntries(a *AppendEntriesArgs, b *AppendEntriesReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e := r.check(); e != nil {
		return e
	}
	if a.LeaderId < 0 || a.LeaderId >= len(r.peers) || a.PrevLogIndex < 0 || a.PrevLogTerm < 0 || a.LeaderCommit < 0 || a.Term < 0 || len(a.Entries) > maxEntries {
		return errors.New("invalid append request")
	}
	b.Success = false
	if a.Term > r.currentTerm {
		if e := r.stepDown(a.Term); e != nil {
			return e
		}
	}
	b.Term = r.currentTerm
	b.NextIndex = len(r.log)
	if a.Term < r.currentTerm {
		return nil
	}
	r.role = Follower
	r.resetDeadline()
	if a.PrevLogIndex >= len(r.log) {
		return nil
	}
	if r.log[a.PrevLogIndex].Term != a.PrevLogTerm {
		i := a.PrevLogIndex
		for i > 0 && r.log[i-1].Term == r.log[a.PrevLogIndex].Term {
			i--
		}
		b.NextIndex = i
		if b.NextIndex < 1 {
			b.NextIndex = 1
		}
		return nil
	}
	if a.PrevLogIndex+len(a.Entries)+1 > maxEntries {
		return errors.New("log capacity reached")
	}
	prev := a.PrevLogTerm
	for _, entry := range a.Entries {
		if entry.Term < prev || entry.Term > a.Term || entry.Term <= 0 || len(entry.Command) > maxCommand {
			return errors.New("invalid log entry")
		}
		prev = entry.Term
	}
	changed := false
	for i, entry := range a.Entries {
		index := a.PrevLogIndex + 1 + i
		if index < len(r.log) {
			if r.log[index].Term == entry.Term {
				if !bytes.Equal(r.log[index].Command, entry.Command) || r.log[index].Noop != entry.Noop {
					return errors.New("conflicting content in same term")
				}
				continue
			}
			if index <= r.commitIndex {
				return errors.New("cannot replace committed entry")
			}
			r.log = r.log[:index]
		}
		r.log = append(r.log, cloneEntries(a.Entries[i:])...)
		changed = true
		break
	}
	matched := a.PrevLogIndex + len(a.Entries)
	commit := a.LeaderCommit
	if commit > matched {
		commit = matched
	}
	if commit > r.commitIndex {
		r.commitIndex = commit
		changed = true
	}
	if changed {
		if e := r.persist(); e != nil {
			return e
		}
	}
	b.Success = true
	b.NextIndex = matched + 1
	return nil
}
func (r *Raft) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started || r.stopped {
		return
	}
	r.started = true
	r.resetDeadline()
	r.wg.Add(1)
	go r.ticker()
}
func (r *Raft) Stop() {
	r.mu.Lock()
	if !r.stopped {
		r.stopped = true
		close(r.stop)
	}
	r.mu.Unlock()
	r.wg.Wait()
}
func (r *Raft) ticker() {
	defer r.wg.Done()
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			r.mu.Lock()
			if r.check() != nil {
				r.mu.Unlock()
				continue
			}
			leader := r.role == Leader
			elect := !leader && time.Now().After(r.deadline)
			send := leader && time.Since(r.lastHeartbeat) >= r.heartbeat
			if send {
				r.lastHeartbeat = time.Now()
			}
			r.mu.Unlock()
			if elect {
				r.startElection()
			}
			if send {
				r.replicate()
			}
		}
	}
}
func (r *Raft) startElection() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.check() != nil {
		return
	}
	r.currentTerm++
	r.votedFor = r.id
	r.role = Candidate
	r.resetDeadline()
	if r.persist() != nil {
		return
	}
	term := r.currentTerm
	last := len(r.log) - 1
	args := RequestVoteArgs{Term: term, CandidateId: r.id, LastLogIndex: last, LastLogTerm: r.log[last].Term}
	votes := 1
	if votes > len(r.peers)/2 {
		r.becomeLeader()
		return
	}
	for i, address := range r.peers {
		if i == r.id {
			continue
		}
		r.wg.Add(1)
		go func(address string) {
			defer r.wg.Done()
			var reply RequestVoteReply
			if !r.sendRPC(address, "Raft.RequestVote", &args, &reply) {
				return
			}
			r.mu.Lock()
			defer r.mu.Unlock()
			if r.check() != nil {
				return
			}
			if reply.Term > r.currentTerm {
				r.stepDown(reply.Term)
				return
			}
			if r.role != Candidate || r.currentTerm != term {
				return
			}
			if reply.VoteGranted {
				votes++
				if votes > len(r.peers)/2 {
					r.becomeLeader()
				}
			}
		}(address)
	}
}
func (r *Raft) becomeLeader() {
	if len(r.log) >= maxEntries {
		r.err = errors.New("log capacity reached")
		return
	}
	r.role = Leader
	r.log = append(r.log, LogEntry{Term: r.currentTerm, Noop: true})
	if r.persist() != nil {
		return
	}
	r.nextIndex = make([]int, len(r.peers))
	r.matchIndex = make([]int, len(r.peers))
	for i := range r.nextIndex {
		r.nextIndex[i] = len(r.log)
	}
	r.matchIndex[r.id] = len(r.log) - 1
	r.advanceCommit()
	r.lastHeartbeat = time.Time{}
}
func (r *Raft) advanceCommit() {
	for n := len(r.log) - 1; n > r.commitIndex; n-- {
		if r.log[n].Term != r.currentTerm {
			continue
		}
		votes := 0
		for _, matched := range r.matchIndex {
			if matched >= n {
				votes++
			}
		}
		if votes > len(r.peers)/2 {
			r.commitIndex = n
			r.persist()
			return
		}
	}
}

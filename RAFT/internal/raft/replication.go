package raft

import (
	"errors"
)

func (r *Raft) replicate() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.check() != nil || r.role != Leader {
		return
	}
	for i, address := range r.peers {
		if i == r.id || r.inflight[i] {
			continue
		}
		r.inflight[i] = true
		next := r.nextIndex[i]
		args := AppendEntriesArgs{Term: r.currentTerm, LeaderId: r.id, PrevLogIndex: next - 1, PrevLogTerm: r.log[next-1].Term, Entries: cloneEntries(r.log[next:]), LeaderCommit: r.commitIndex}
		r.wg.Add(1)
		go func(i int, address string, args AppendEntriesArgs, next int) {
			defer r.wg.Done()
			var reply AppendEntriesReply
			ok := r.sendRPC(address, "Raft.AppendEntries", &args, &reply)
			r.mu.Lock()
			defer r.mu.Unlock()
			r.inflight[i] = false
			if !ok || r.check() != nil {
				return
			}
			if reply.Term > r.currentTerm {
				r.stepDown(reply.Term)
				return
			}
			if r.role != Leader || r.currentTerm != args.Term {
				return
			}
			if reply.Success {
				matched := args.PrevLogIndex + len(args.Entries)
				if matched > r.matchIndex[i] {
					r.matchIndex[i] = matched
					r.nextIndex[i] = matched + 1
				}
				r.advanceCommit()
			} else if r.nextIndex[i] == next {
				n := reply.NextIndex
				if n < 1 {
					n = 1
				}
				if n >= next {
					n = next - 1
				}
				if n < 1 {
					n = 1
				}
				r.nextIndex[i] = n
			}
		}(i, address, args, next)
	}
}

// Propose durably appends on the leader; replication/commit is asynchronous.
func (r *Raft) Propose(command []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e := r.check(); e != nil {
		return 0, e
	}
	if r.role != Leader {
		return 0, errors.New("not leader")
	}
	if len(command) > maxCommand || len(r.log) >= maxEntries {
		return 0, errors.New("log capacity exceeded")
	}
	r.log = append(r.log, LogEntry{Term: r.currentTerm, Command: append([]byte(nil), command...)})
	if e := r.persist(); e != nil {
		return 0, e
	}
	index := len(r.log) - 1
	r.matchIndex[r.id] = index
	r.advanceCommit()
	if r.err != nil {
		return 0, r.err
	}
	return index, nil
}
func (r *Raft) Submit(a *SubmitArgs, b *SubmitReply) error {
	index, e := r.Propose(a.Command)
	r.mu.Lock()
	defer r.mu.Unlock()
	b.Index = index
	b.Term = r.currentTerm
	b.Leader = r.role == Leader
	return e
}
func (r *Raft) Status(a *StatusArgs, b *StatusReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e := r.check(); e != nil {
		return e
	}
	b.Term = r.currentTerm
	b.CommitIndex = r.commitIndex
	b.LastLogIndex = len(r.log) - 1
	b.Role = r.role
	return nil
}

// Applied is an ordered snapshot of committed application commands. Consumers
// retain their own applied index; no-op leadership entries are omitted.
func (r *Raft) Applied() []AppliedEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []AppliedEntry{}
	for i := 1; i <= r.commitIndex; i++ {
		if !r.log[i].Noop {
			out = append(out, AppliedEntry{i, append([]byte(nil), r.log[i].Command...)})
		}
	}
	return out
}

package raft

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

const maxCommand = 1 << 20
const maxEntries = 10000
const maxStateBytes = 64 << 20

type diskState struct {
	Version, ID, Term, Vote, Commit int
	Peers                           []string
	Log                             []LogEntry
}

func NewPersistentRaft(id int, peers []string, send func(string, string, interface{}, interface{}) bool, path string) (*Raft, error) {
	if len(peers) == 0 || id < 0 || id >= len(peers) || send == nil {
		return nil, errors.New("invalid Raft configuration")
	}
	seen := map[string]bool{}
	for _, p := range peers {
		if p == "" || seen[p] {
			return nil, errors.New("empty or duplicate peer")
		}
		seen[p] = true
	}
	r := &Raft{id: id, peers: append([]string(nil), peers...), votedFor: -1, log: []LogEntry{{}}, sendRPC: send, path: path, stop: make(chan struct{}), heartbeat: 80 * time.Millisecond, election: 600 * time.Millisecond, inflight: make([]bool, len(peers))}
	if path != "" {
		if stat, e := os.Stat(path); e == nil && stat.Size() > maxStateBytes {
			return nil, errors.New("persistent state exceeds size limit")
		}
		b, e := os.ReadFile(path)
		if e == nil {
			var s diskState
			if e = json.Unmarshal(b, &s); e != nil {
				return nil, e
			}
			if s.Version != 1 || s.ID != id || !reflect.DeepEqual(s.Peers, peers) || s.Term < 0 || s.Vote < -1 || s.Vote >= len(peers) || len(s.Log) == 0 || len(s.Log) > maxEntries || s.Commit < 0 || s.Commit >= len(s.Log) || s.Log[0].Term != 0 {
				return nil, errors.New("invalid persistent Raft state")
			}
			total := 0
			prev := 0
			for i, entry := range s.Log {
				if entry.Term < prev || entry.Term > s.Term || len(entry.Command) > maxCommand || i > 0 && entry.Term == 0 {
					return nil, errors.New("invalid persistent log")
				}
				total += len(entry.Command)
				if total > maxStateBytes/2 {
					return nil, errors.New("log payload limit exceeded")
				}
				prev = entry.Term
			}
			r.currentTerm = s.Term
			r.votedFor = s.Vote
			r.log = s.Log
			r.commitIndex = s.Commit
		} else if !os.IsNotExist(e) {
			return nil, e
		} else {
			if e = r.persist(); e != nil {
				return nil, e
			}
		}
	}
	r.resetDeadline()
	return r, nil
}
func (r *Raft) persist() error {
	if r.path == "" {
		return nil
	}
	dir := filepath.Dir(r.path)
	if e := os.MkdirAll(dir, 0700); e != nil {
		r.err = e
		return e
	}
	s := diskState{Version: 1, ID: r.id, Term: r.currentTerm, Vote: r.votedFor, Commit: r.commitIndex, Peers: r.peers, Log: r.log}
	b, e := json.Marshal(s)
	if e != nil {
		r.err = e
		return e
	}
	f, e := os.CreateTemp(dir, ".raft-*")
	if e != nil {
		r.err = e
		return e
	}
	name := f.Name()
	defer os.Remove(name)
	if _, e = f.Write(b); e == nil {
		e = f.Sync()
	}
	closeErr := f.Close()
	if e == nil {
		e = closeErr
	}
	if e == nil {
		e = os.Rename(name, r.path)
	}
	if e == nil {
		var d *os.File
		d, e = os.Open(dir)
		if e == nil {
			e = d.Sync()
			d.Close()
		}
	}
	if e != nil {
		r.err = e
	}
	return e
}
func (r *Raft) check() error {
	if r.stopped {
		return errors.New("Raft stopped")
	}
	return r.err
}
func cloneEntries(entries []LogEntry) []LogEntry {
	out := make([]LogEntry, len(entries))
	for i, e := range entries {
		out[i] = e
		out[i].Command = append([]byte(nil), e.Command...)
	}
	return out
}

func (r *Raft) payloadSize() int {
	total := 0
	for _, e := range r.log {
		total += len(e.Command)
	}
	return total
}

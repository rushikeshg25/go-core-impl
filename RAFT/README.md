# Persistent Raft v1

An educational static-membership Raft implementation with election, durable votes/logs, AppendEntries consistency and conflict repair, majority commit and ordered application checkpoints.

Start three terminals from this directory:

```sh
go run ./cmd/raft-demo -id 0
go run ./cmd/raft-demo -id 1
go run ./cmd/raft-demo -id 2
```

The default ordered cluster is localhost:8000,8001,8002. Nodes store state separately under `data/node-<id>/state.json`. Use `-cluster` and `-state` to configure addresses and paths; the persisted identity/peer list must match on restart.

```sh
go run ./cmd/raft-demo -status localhost:8000
# role 2 is leader; submit to whichever node currently leads:
go run ./cmd/raft-demo -submit localhost:8000 -command hello
```

Submission returns a durably appended index, not a commit acknowledgement. A majority must replicate it before application. A timed-out RPC has an unknown outcome; v1 has no client request deduplication.

## Safety and application contract

Candidates must have an up-to-date log. Term/vote/log state is atomically replaced and fsynced, including the directory, before successful acknowledgements. A persistence error stops subsequent protocol operations. Leaders append a no-op for their term and advance commit only when a majority stores a current-term entry. Followers reject conflicting committed entries and repair uncommitted suffixes. Committed indexes are persisted so restart exposes the committed prefix.

`Propose([]byte)` copies and appends on the leader. `Applied()` returns an owned ordered snapshot of committed application commands, omitting leadership no-ops; it returns nil after stop or persistence failure (use `ApplyTo` for an explicit health error). `ApplyTo(lastApplied, callback)` applies a captured committed prefix in order outside the Raft lock and returns the last successfully applied index; callback failure leaves that command unacknowledged. The caller serializes application and durably checkpoints the index with its state machine. This interface does not promise exactly-once external side effects.

`Raft.Start` and `Raft.Stop` are idempotent. RPC dialing/calls are bounded; server shutdown closes active connections. The legacy `NewRaft` constructor is in-memory only; the demo uses `NewPersistentRaft`.

## Boundaries

V1 assumes trusted fixed peers and exclusive ownership of each state file. It has no TLS/authentication, membership changes, snapshots, linearizable reads or production deployment guarantee. State is rewritten as an atomic JSON snapshot on changes; limits are 10,000 log entries including the sentinel, 1 MiB per command and 32 MiB total command payload. Log compaction is deferred; reaching a limit rejects new work. Existing election-only state had no on-disk format to migrate. Creation of a new state directory does not fsync that directory's own parent.

Run `go test -race ./...`. Tests cover minority isolation, leader replacement, conflicting suffix repair, durable votes, stale-candidate rejection, restart, persistence failures, application callback retry and a real three-node TCP/RPC cluster. These are bounded acceptance tests, not a proof of consensus correctness under every failure schedule.

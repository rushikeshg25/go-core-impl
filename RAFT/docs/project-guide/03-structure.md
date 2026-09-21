# Structure

## What lives where

This small Go module separates process wiring in `cmd/raft-demo`, consensus and durable state in `internal/raft`, and the TCP adapter in `internal/rpc`. Both implementation packages use Go's `internal` import boundary; the demo is their in-module consumer ([go.mod:1](../../go.mod#L1), [main.go:6](../../cmd/raft-demo/main.go#L6)). The inspected scope contains 14 tracked files before this guide, including eight Go source/test files; there is no local CI, deployment manifest, or third-party dependency list.

```text
RAFT/
├── cmd/raft-demo/main.go      node and client entry point
├── internal/
│   ├── raft/                 core, persistence, replication, tests
│   ├── rpc/                  TCP transport and integration test
│   └── storage/.gitkeep       empty placeholder; not the persistence implementation
├── go.mod                    module and Go version
├── README.md                 usage and current boundaries
├── V1.md                     delivery contract
├── HISTORY.md                delivery chronology and prior validation
└── docs/project-guide/       this guide
```

## Demo

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [cmd/raft-demo/main.go](../../cmd/raft-demo/main.go#L15) | Flag parsing; status/submit clients; node initialization; signal shutdown. | `main` executable entry point; no exported API | `go run ./cmd/raft-demo` |

## Internal raft

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [internal/raft/types.go](../../internal/raft/types.go#L8) | Role constants, core state, wire/request/reply types, application entry, legacy in-memory constructor. | `Raft`, `NodeRole`, `LogEntry`, `RequestVoteArgs/Reply`, `AppendEntriesArgs/Reply`, `SubmitArgs/Reply`, `StatusArgs/Reply`, `AppliedEntry`, `NewRaft` | All core files; demo, transport, tests |
| [internal/raft/storage.go](../../internal/raft/storage.go#L22) | Configuration/snapshot validation; load/init; JSON atomic replacement; latched errors; copy and size helpers. | `NewPersistentRaft` | Demo/test setup; node and replication mutations |
| [internal/raft/node.go](../../internal/raft/node.go#L22) | Vote and append receivers; lifecycle/ticker; elections; no-op leadership entry; majority commit. | `RequestVote`, `AppendEntries`, `Start`, `Stop` | RPC dispatcher, demo, core ticker, tests |
| [internal/raft/replication.go](../../internal/raft/replication.go#L7) | Per-peer suffix replication; proposals; client-facing RPC handlers; copied application snapshots/checkpoints. | `Propose`, `Submit`, `Status`, `Applied`, `ApplyTo` | Ticker, RPC dispatcher, application callers, tests |
| [internal/raft/raft_test.go](../../internal/raft/raft_test.go#L11) | In-process injectable network; partition/repair/restart; durable vote/freshness; single-node/persistence failure; ordered callback retry. | `TestPartitionRepairAndRestart`, `TestVoteDurabilityAndFreshness`, `TestSingleNodeAndPersistenceFailure`, `TestOrderedApplicationCheckpoint`, `TestFailedCommitPersistenceHidesApplicationSnapshot` | Go test runner; same-package access to protocol internals |

The storage implementation is [storage.go](../../internal/raft/storage.go#L74), despite the separate empty `internal/storage` directory. Tests manually trigger elections/replication to control scenarios; the TCP test exercises the actual ticker ([raft_test.go:62](../../internal/raft/raft_test.go#L62), [rpc_test.go:34](../../internal/rpc/rpc_test.go#L34)).

## Internal rpc

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [internal/rpc/rpc.go](../../internal/rpc/rpc.go#L12) | Register node methods; accept/track TCP connections; stop/wait; bounded outgoing RPC. | `Server`, `NewServer`, `Start`, `Serve`, `Stop`, `Call` | Demo; injected sender; transport test |
| [internal/rpc/rpc_test.go](../../internal/rpc/rpc_test.go#L11) | Three real loopback listeners, persistent nodes, leader discovery, submit, eventual copied application result on every node. | `TestRealClusterElectionAndReplication` | Go test runner |

## Module root

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [go.mod](../../go.mod#L1) | Module identity and Go version; no third-party requirements. | Module `github.com/rushikeshg25/raft` | Go toolchain |
| [README.md](../../README.md#L1) | Run commands, API contract, limits, trust assumptions. | User-facing documentation | Maintainers/operators |
| [V1.md](../../V1.md#L3) | v1 feature and acceptance scope. | Delivery contract | Maintainers/reviewers |
| [HISTORY.md](../../HISTORY.md#L10) | Referenced implementation milestones and past verification. | Evidence chronology | Maintainers/reviewers |

## Excluded

[.gitignore](../../.gitignore) is repository housekeeping. [internal/storage/.gitkeep](../../internal/storage/.gitkeep) is an empty placeholder, not a source module. Runtime `data/` files are generated node state, never source; their path comes from [main.go:40](../../cmd/raft-demo/main.go#L40). This guide excludes the parent monorepo and its sibling projects.

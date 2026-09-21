# Decisions

This document separates choices stated in the contract/comments from rationales inferred from implementation. The delivery history records milestones; commit titles alone are not evidence of correctness ([HISTORY.md:12](../../HISTORY.md#L12)).

## Fixed membership is part of durable identity

- **What:** Ordered peer addresses and the node ID are stored alongside the consensus state and must match exactly at restart.
- **Evidence:** [storage.go:16](../../internal/raft/storage.go#L16), [storage.go:44](../../internal/raft/storage.go#L44); membership changes are excluded by [V1.md:5](../../V1.md#L5).
- **Why:** Explicit v1 scope keeps configuration changes outside this implementation.
- **Tradeoff:** Fewer consensus transitions to reason about, but address/order changes cannot be made by simply editing flags against existing state.
- **Confidence:** Scope confirmed by the contract; simplification rationale inferred.

## Whole-state atomic snapshots favor a small implementation

- **What:** Persist the complete JSON state with a same-directory temporary file, fsync, rename, and directory fsync. A persistence error latches into the node.
- **Evidence:** [storage.go:74](../../internal/raft/storage.go#L74), [storage.go:119](../../internal/raft/storage.go#L119); persistent votes and failure rejection have scenario coverage in [raft_test.go:113](../../internal/raft/raft_test.go#L113) and [raft_test.go:147](../../internal/raft/raft_test.go#L147).
- **Why, inferred:** One format and replacement procedure keep durability understandable without a segmented write-ahead log.
- **Tradeoff:** Every update serializes/writes the accumulated state under the core mutex; failures leave the object unusable for guarded operations rather than automatically retrying. Directory creation does not fsync the directory's parent, an explicit durability boundary ([README.md:33](../../README.md#L33)).
- **Confidence:** Mechanism confirmed; rationale inferred.

## Current-term no-ops support majority commit after leadership changes

- **What:** Every new leader appends a no-op in its own term. Commit advancement selects only current-term indexes held by a strict majority, making all earlier indexes committed with that prefix.
- **Evidence:** [node.go:247](../../internal/raft/node.go#L247), [node.go:266](../../internal/raft/node.go#L266); application filters no-ops at [replication.go:119](../../internal/raft/replication.go#L119).
- **Why, inferred:** A leadership entry gives the new term something to replicate even when no client command arrives.
- **Tradeoff:** Leadership consumes log capacity and application indexes contain gaps when viewed through `Applied()`.
- **Confidence:** Behavior confirmed by source and [README.md:25](../../README.md#L25); rationale inferred.

## An append acknowledgement is deliberately separate from application

- **What:** `Propose` persists locally and returns an index, while periodic replication establishes commit later. `ApplyTo` delegates effects and checkpoint durability to its caller.
- **Evidence:** Explicit comments at [replication.go:62](../../internal/raft/replication.go#L62) and [replication.go:127](../../internal/raft/replication.go#L127); command bytes contain no client/request identity ([types.go:55](../../internal/raft/types.go#L55)).
- **Why:** The API explicitly makes replication asynchronous and leaves the state machine external.
- **Tradeoff:** Clients must distinguish append from commitment; retries can duplicate commands. A callback that performs an effect and then fails can repeat that effect on retry. External exactly-once effects are not promised ([README.md:21](../../README.md#L21), [README.md:27](../../README.md#L27)).
- **Confidence:** Confirmed by comments and documented contract.

## Copy at ownership boundaries and release the lock for application callbacks

- **What:** Copy proposal inputs, replicated suffixes, and application results. Capture committed entries under the mutex, then invoke callbacks after unlocking.
- **Evidence:** [storage.go:125](../../internal/raft/storage.go#L125), [replication.go:75](../../internal/raft/replication.go#L75), [replication.go:140](../../internal/raft/replication.go#L140). Tests mutate input bytes and reenter `Status` from an application callback ([raft_test.go:68](../../internal/raft/raft_test.go#L68), [raft_test.go:173](../../internal/raft/raft_test.go#L173)).
- **Why, inferred:** Avoid caller mutation of consensus state and avoid holding the protocol mutex while arbitrary application code runs.
- **Tradeoff:** Copy cost grows with log size; concurrent applications still require caller serialization, and each call sees only its captured prefix.
- **Confidence:** Behavior/tests confirmed; rationale inferred.

## Transport stays injectable and intentionally small

- **What:** The core accepts one generic RPC function. The supplied transport opens a fresh connection per call; replication allows one active request per peer and retries on later heartbeats.
- **Evidence:** [types.go:33](../../internal/raft/types.go#L33), [rpc.go:81](../../internal/rpc/rpc.go#L81), [replication.go:13](../../internal/raft/replication.go#L13); tests substitute [raft_test.go:17](../../internal/raft/raft_test.go#L17).
- **Why, inferred:** The same core can be exercised in controlled partitions and over real TCP with little adapter code.
- **Tradeoff:** Boolean results erase detailed errors, connections are not pooled, and a slow or permanently blocked custom sender undermines bounded shutdown. TLS/authentication and hostile-network operation are outside the stated scope ([README.md:33](../../README.md#L33)).
- **Confidence:** Mechanism confirmed; rationale inferred.

## Gotchas

- **Node and server lifecycle contracts differ.** [README.md:29](../../README.md#L29) scopes idempotence to `Raft.Start` and `Raft.Stop`; `Raft.Start` ignores repeats, but `Server.Serve` rejects an already-started/stopped server, and `Server.Start` attempts its bind first ([node.go:144](../../internal/raft/node.go#L144), [rpc.go:22](../../internal/rpc/rpc.go#L22)). The transport and node lifecycle methods have different repeat-call behavior.
- **`NewRaft` is not durable.** It calls the persistent constructor with an empty path, and panics on configuration errors; the demo uses the error-returning durable constructor ([types.go:70](../../internal/raft/types.go#L70), [main.go:43](../../cmd/raft-demo/main.go#L43)).
- **`Applied()` hides snapshots after stop or persistence failure.** It returns nil when unhealthy; `ApplyTo` and `Status` return explicit health errors ([replication.go:112](../../internal/raft/replication.go#L112), [replication.go:130](../../internal/raft/replication.go#L130)). Do not infer recovery or durability from that accessor alone.
- **No background application exists.** The demo accepts opaque bytes but never calls `ApplyTo`; a state machine is integration work for a consumer ([main.go:43](../../cmd/raft-demo/main.go#L43), [replication.go:127](../../internal/raft/replication.go#L127)).
- **Log size is a lifetime constraint.** The sentinel and every leadership no-op count toward 10,000 entries; capacity exhaustion can prevent leadership, not just reject proposals ([storage.go:13](../../internal/raft/storage.go#L13), [node.go:247](../../internal/raft/node.go#L247)).
- **Status does not certify a read.** It returns mutex-protected local fields without contacting a quorum; an isolated leader can still report its role until it learns a higher term ([replication.go:96](../../internal/raft/replication.go#L96), [replication.go:31](../../internal/raft/replication.go#L31)).
- **Tests establish bounded scenarios.** Five core tests and one TCP test cover meaningful paths, but the broad invalid-input/lifecycle acceptance wording is not an exhaustive current test inventory ([V1.md:9](../../V1.md#L9), [raft_test.go](../../internal/raft/raft_test.go), [rpc_test.go](../../internal/rpc/rpc_test.go)).

## Conventions

Core mutations generally run under `r.mu`; unexported helpers such as `persist`, `stepDown`, and `advanceCommit` rely on their callers for locking ([node.go:22](../../internal/raft/node.go#L22), [replication.go:63](../../internal/raft/replication.go#L63)). RPC handlers use pointer argument/reply structs and return `error`, while outbound transport returns `bool` ([types.go:40](../../internal/raft/types.go#L40), [rpc.go:81](../../internal/rpc/rpc.go#L81)). Tests live beside implementation and use the same package, intentionally allowing direct access to internals for scenario control ([raft_test.go:1](../../internal/raft/raft_test.go#L1), [raft_test.go:62](../../internal/raft/raft_test.go#L62)).

# Persistent static-membership Raft history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2026-02-11T22:03:38+05:30: feat:RAFT impl

- **What happened:** The repository records `feat:RAFT impl`.
- **Evidence:** [commit 1f4ff2d431](https://github.com/rushikeshg25/go-core-impl/commit/1f4ff2d4313dd795c9c2ac016715f41dd8392a1f).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:21:38+05:30: docs: define raft v1 contract

- **What happened:** The repository records `docs: define raft v1 contract`.
- **Evidence:** [commit 59e3d5bccc](https://github.com/rushikeshg25/go-core-impl/commit/59e3d5bccc940456bf39e911424e8f5fd1386f00).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:54:12+05:30: feat(raft): persist node identity term vote log and committed prefix atomically

- **What happened:** The repository records `feat(raft): persist node identity term vote log and committed prefix atomically`.
- **Evidence:** [commit 23200edbec](https://github.com/rushikeshg25/go-core-impl/commit/23200edbec08f2a693c5330d90aaa125502ce86b).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:54:12+05:30: feat(raft): enforce election freshness and durable AppendEntries consistency

- **What happened:** The repository records `feat(raft): enforce election freshness and durable AppendEntries consistency`.
- **Evidence:** [commit 0cf3ed1812](https://github.com/rushikeshg25/go-core-impl/commit/0cf3ed181206f74d7cbd153f18d651c2fa7e19e9).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:54:12+05:30: feat(raft): replicate repair conflicting suffixes and commit by majority

- **What happened:** The repository records `feat(raft): replicate repair conflicting suffixes and commit by majority`.
- **Evidence:** [commit 3e3221c8e7](https://github.com/rushikeshg25/go-core-impl/commit/3e3221c8e7c6be0e8d0a040939b2873be2fb87fa).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:56:34+05:30: test(raft): verify majority safety partition repair durable votes and restart

- **What happened:** The repository records `test(raft): verify majority safety partition repair durable votes and restart`.
- **Evidence:** [commit fdd8544385](https://github.com/rushikeshg25/go-core-impl/commit/fdd85443858a91da56cb1999d85a39dc5fcff526).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:56:34+05:30: feat(raft): add bounded RPC lifecycle persistent demo and real cluster verification

- **What happened:** The repository records `feat(raft): add bounded RPC lifecycle persistent demo and real cluster verification`.
- **Evidence:** [commit 18923912e3](https://github.com/rushikeshg25/go-core-impl/commit/18923912e303d05d6ec795cb0c3874469a03a21f).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:59:21+05:30: fix(raft): bound persisted state and avoid stale leader election triggers

- **What happened:** The repository records `fix(raft): bound persisted state and avoid stale leader election triggers`.
- **Evidence:** [commit 7c2a071b98](https://github.com/rushikeshg25/go-core-impl/commit/7c2a071b9880581ed886a3303265e51aae09a025).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T14:01:39+05:30: feat(raft): apply committed commands through resumable ordered checkpoints

- **What happened:** The repository records `feat(raft): apply committed commands through resumable ordered checkpoints`.
- **Evidence:** [commit 9c802ac263](https://github.com/rushikeshg25/go-core-impl/commit/9c802ac26336fd744addd3086821b6ed7934bbd2).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T14:04:58+05:30: docs: document raft v1 usage and limitations

- **What happened:** The repository records `docs: document raft v1 usage and limitations`.
- **Evidence:** [commit 21f5ef061a](https://github.com/rushikeshg25/go-core-impl/commit/21f5ef061a2f24c8a57a4a9e607ad909cd000fc5).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed `go test -race ./...`. [raft_test.go](internal/raft/raft_test.go) exercises deterministic partitions, voting/restart and application checkpoints; [rpc_test.go](internal/rpc/rpc_test.go) verifies election and replication over a real three-node TCP cluster.

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.

## 2026-09-21: Project guide and failure snapshot follow-up

Added the six-file [project guide](docs/project-guide/README.md), traced against the v1 source. The guide read-through identified that `Applied` could expose in-memory state after commit persistence failed. It now returns nil on a stopped or failed node, while `ApplyTo` returns the explicit error. [The regression test](internal/raft/raft_test.go) forces commit persistence failure and checks both APIs. The full `go test -race ./...` suite, including real TCP replication, passed again after this fix.

# RAFT Project Guide

> Generated: 2026-09-21 from the v1 checkout, then refreshed for the failure-snapshot guard in this delivery. Scope: the RAFT module only. Source and tests were read; the primary verification reran `go test -race ./...` successfully after the guard was added.

## What this is

RAFT is an educational implementation of fixed-membership consensus with durable terms, votes, logs, and a committed prefix ([v1 contract](../../V1.md#contract)). Its audience is a developer studying election, replication, recovery, and ordered application rather than operating a production consensus service ([boundaries](../../README.md#boundaries)). It consists of an internal Go library, a TCP/RPC adapter, and a demo executable with node, status, and submission modes ([main.go:15](../../cmd/raft-demo/main.go#L15)).

## The five-file tour

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [cmd/raft-demo/main.go](../../cmd/raft-demo/main.go#L15) | Connects flags, durable nodes, transport, and shutdown. | [Startup](02-flow.md#startup) |
| 2 | [internal/raft/types.go](../../internal/raft/types.go#L8) | Names roles, log entries, node state, and wire messages. | [Data model](01-architecture.md#data-model) |
| 3 | [internal/raft/storage.go](../../internal/raft/storage.go#L22) | Establishes durable identity and the failure boundary. | [State and persistence](01-architecture.md#state-and-persistence) |
| 4 | [internal/raft/node.go](../../internal/raft/node.go#L22) | Explains votes, follower consistency, election, and majority commit. | [Election](02-flow.md#election) |
| 5 | [internal/raft/replication.go](../../internal/raft/replication.go#L7) | Completes proposal, replication, status, and application APIs. | [Submission and replication](02-flow.md#submission-and-replication) |

## Run it

Use Go 1.24.6 as declared in [go.mod:3](../../go.mod#L3). Run each of these in a separate terminal from the RAFT module directory:

```sh
go run ./cmd/raft-demo -id 0
go run ./cmd/raft-demo -id 1
go run ./cmd/raft-demo -id 2
```

The default ordered peer list is `localhost:8000,localhost:8001,localhost:8002`; each node writes `data/node-<id>/state.json`. Override with `-cluster` and `-state`, keeping the same ID and ordered peer list when reopening existing state ([main.go:16](../../cmd/raft-demo/main.go#L16), [storage.go:44](../../internal/raft/storage.go#L44)). Configuration is through flags; this entry point reads no environment variables.

```sh
go run ./cmd/raft-demo -status localhost:8000
go run ./cmd/raft-demo -status localhost:8001
go run ./cmd/raft-demo -status localhost:8002
# Replace the address below with the node reporting role=2.
go run ./cmd/raft-demo -submit localhost:8000 -command hello
# Project verification command:
go test -race ./...
```

Role values are follower `0`, candidate `1`, leader `2` ([types.go:10](../../internal/raft/types.go#L10)). Submission acknowledges a durable local append; replication and commit follow asynchronously, and the demo does not execute an application callback ([replication.go:62](../../internal/raft/replication.go#L62), [main.go:31](../../cmd/raft-demo/main.go#L31)). A local status result is diagnostic rather than a linearizable read or a deduplicated client completion protocol ([boundaries](../../README.md#boundaries)). [HISTORY.md:74](../../HISTORY.md#L74) records an earlier successful race-test run; that is historical evidence, not a run performed for this guide.

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, contracts, durability, and deployment.
2. [Flow](02-flow.md) — startup through append, majority commit, application, and shutdown.
3. [Structure](03-structure.md) — every significant source file and its callers.
4. [Tech stack](04-tech-stack.md) — Go, standard-library transport/storage, and test tooling.
5. [Decisions](05-decisions.md) — explicit boundaries, inferred tradeoffs, and traps.

## Open questions

- Which application state machine will consume `ApplyTo`, and how will it atomically checkpoint effects with the returned index? The library assigns that responsibility to the caller; the demo has no consumer ([replication.go:127](../../internal/raft/replication.go#L127), [main.go:51](../../cmd/raft-demo/main.go#L51)).
- What is the operational response when the bounded log fills? v1 excludes snapshots and membership changes, and leadership itself needs a free no-op slot ([V1.md:5](../../V1.md#L5), [node.go:247](../../internal/raft/node.go#L247)).
- How will invalid-input and lifecycle acceptance be extended? [V1.md:9](../../V1.md#L9) requests such tests; the current six tests concentrate on partitions, persistence, application retry, and real TCP replication ([test map](03-structure.md#internal-raft)). No production or long-running operational validation is recorded ([HISTORY.md:82](../../HISTORY.md#L82)).

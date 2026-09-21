# Durable logs project guide

> Generated: 2026-09-21 from commit `f89ecbf`. Scope: this `durable-logs/` Go module.

## What this is

This library appends timestamped messages to numbered local files using length-prefixed protobuf records ([writer](../../durablelogs/core.go#L80)). Go callers use an explicitly flushed, mutex-protected logger that validates existing segments and can repair an incomplete final tail on reopen ([Open](../../durablelogs/core.go#L30)). A small executable demonstrates ten writes and close, with no network service or background worker ([main](../../main.go#L8)).

## The five-file tour

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [main.go](../../main.go#L8) | Open, append, and close at the call site. | [Startup](02-flow.md#startup-and-recovery) |
| 2 | [durablelogs/core.go](../../durablelogs/core.go#L18) | Owns buffering, synchronization, rotation, and lifecycle. | [Components](01-architecture.md#components) |
| 3 | [log.proto](../../log.proto#L7) | The two stored fields. | [Data model](01-architecture.md#data-model) |
| 4 | [durablelogs/records.go](../../durablelogs/records.go#L39) | Framing validation and narrowly scoped repair. | [Recovery](02-flow.md#startup-and-recovery) |
| 5 | [durablelogs/core_test.go](../../durablelogs/core_test.go#L10) | Executable examples of restart, concurrency, and failures. | [Gotchas](05-decisions.md#gotchas) |

## Run it

From this module directory, with Go supporting the [1.23.5 directive](../../go.mod#L3):

```sh
go mod download
go run .
go test -race ./...
```

The demo writes ten `hello` messages under `./logs`, rotating at five records per segment; reruns append to existing data ([main.go:9](../../main.go#L9)). There are no required environment variables in the entry point. The existing delivery records a passing race test run in [HISTORY.md](../../HISTORY.md#delivery-verification); this guide does not constitute another test run.

## Reading order for this guide

1. [Architecture](01-architecture.md) — components and durability boundaries.
2. [Flow](02-flow.md) — startup, append, flush, replay, and close.
3. [Structure](03-structure.md) — complete handwritten source inventory.
4. [Tech stack](04-tech-stack.md) — runtime, protobuf, and tooling evidence.
5. [Decisions](05-decisions.md) — rationale, tradeoffs, and traps.

## Open questions

- Is regeneration of the unused generated `Mesage` type intended? Its [header](../../durablelogs/pb/message.pb.go#L5) references `message.proto`, which is absent from this module. The checked-in [log.proto](../../log.proto#L5) also declares a Go package path different from the [actual import](../../durablelogs/core.go#L5); no generation command is supplied here.
- What workload size and crash guarantees are required? The [README](../../README.md#L21) excludes directory-entry durability, and [ReadAll](../../durablelogs/core.go#L178) materializes the full log. No benchmark or deployment configuration is present in the [scoped inventory](03-structure.md).
- Should validation coverage expand beyond the four [current tests](../../durablelogs/core_test.go#L10)? Invalid segment names, earlier-segment truncation, oversized lengths, invalid protobuf, and input limits have implementation branches but no dedicated tests in that file.

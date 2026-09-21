# Event Loop Project Guide

> Generated: 2026-09-21 from commit `f3bcdc0`. Scope: the `event-loop` module only.

## What this is

This Go demo runs submitted functions with callbacks, with at most five asynchronous tasks outstanding ([dispatcher](../../main.go#L60)). It is a concurrency learning example for developers, packaged as an executable rather than an importable library ([package](../../main.go#L1), [demo](../../main.go#L136)). A coordinator serializes synchronous tasks and callbacks, while shutdown stops admission and drains accepted work ([lifecycle](../../main.go#L97)).

## The five-file tour

There are only two Go source files; the last three stops supply the module and contract context.

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [main.go](../../main.go#L136) | Demo, public API, then dispatcher in one file | [Flow](02-flow.md#startup) |
| 2 | [main_test.go](../../main_test.go#L11) | Observable lifecycle guarantees | [Structure](03-structure.md#module-root) |
| 3 | [go.mod](../../go.mod#L1) | Independent module and Go version | [Stack](04-tech-stack.md#languages-and-runtimes) |
| 4 | [README.md](../../README.md#L5) | Supported usage and blocking constraints | [Decisions](05-decisions.md#gotchas) |
| 5 | [V1.md](../../V1.md#L3) | Intended scope, including cooperative cancellation | [Open questions](#open-questions) |

## Run it

From the repository root:

```bash
cd event-loop
go run .
go test -race ./...
```

Use Go 1.23.5 or a compatible newer toolchain ([go.mod](../../go.mod#L3)); the demo needs no services, arguments, or environment variables ([main](../../main.go#L136)). It prints `task`, then `done`. The [module README](../../README.md#L3) documents these commands; this guide did not rerun tests.

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, contracts, and limits.
2. [Flow](02-flow.md) — startup, dispatch, callback submission, and shutdown.
3. [Structure](03-structure.md) — complete source inventory.
4. [Tech stack](04-tech-stack.md) — runtime and tooling.
5. [Decisions](05-decisions.md) — evidenced choices and traps.

## Open questions

- Should tests cover concurrent stop callers, panic behavior, and nonreturning tasks? The [suite](../../main_test.go#L11) covers drain/bounds, concurrent submission, order, cancellation, and invalid admission, but these failure cases remain untested.
- Is `AddCallback` intentionally allowed to execute a task as well as its callback and ignore `Async`? That is the [current dispatch behavior](../../main.go#L87).

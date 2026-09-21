# WebSocket echo server project guide

> Generated: 2026-09-21 from checkout commit `19d7721`. Scope: the seven tracked files in the `websockets/` Go module.

## What this is

This is a small WebSocket echo server that implements its own upgrade and frame handling using Go's standard library ([handler](../../main.go#L34)). It is an educational implementation for readers exploring masked frames, fragmented messages and control frames ([v1 contract](../../V1.md#L5)). One HTTP process exposes `/ws` and `/health` on port 8080, with no datastore ([startup](../../main.go#L215)).

## Run it

From this module directory, using Go 1.23.2 or a compatible newer toolchain ([manifest](../../go.mod#L3)):

```bash
go run .
# In another terminal:
curl http://localhost:8080/health
# Connect a WebSocket client to ws://localhost:8080/ws
go test -race ./...
go test -run '^$' -fuzz FuzzFrame -fuzztime 5s
```

The address is hardcoded and no environment variables are read ([main](../../main.go#L216)). [ws.sh](../../ws.sh#L6) sends an HTTP upgrade with curl; it does not encode masked WebSocket messages. The checkout's [root makefile](../../../makefile#L9) invokes `make go run .` for this module despite there being no module Makefile; use `go run .` directly.

## The five-file tour

There are only four code, script and manifest files; the fifth stop is the explicit scope document.

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [go.mod](../../go.mod#L1) | Establishes the toolchain and absence of third-party requirements. | [Tech stack](04-tech-stack.md#languages-and-runtimes) |
| 2 | [main.go](../../main.go#L216) | Follow startup into upgrade, message assembly and wire framing. | [Flow](02-flow.md#startup) |
| 3 | [ws.sh](../../ws.sh#L6) | Shows the upgrade headers a client sends. | [Structure](03-structure.md#module-root) |
| 4 | [main_test.go](../../main_test.go#L28) | Shows a real TCP exchange with a ping between fragments. | [Verification flow](02-flow.md#verification-flow) |
| 5 | [V1.md](../../V1.md#L3) | Defines the intended contract and acceptance criteria. | [Decisions](05-decisions.md) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, state and operational boundaries.
2. [Flow](02-flow.md) — startup, upgrade, echo, control and verification paths.
3. [Structure](03-structure.md) — every significant module file.
4. [Tech stack](04-tech-stack.md) — standard library, runtime and tooling.
5. [Decisions](05-decisions.md) — evidenced choices and maintenance traps.

## Open questions

- Is broader protocol conformance testing intended? Current [tests](../../main_test.go#L28) do not exercise binary echo, large length encodings, fragmented size overflow, UTF-8 rejection, origin/version failures or timeout behavior.
- How should production connection limits, TLS and graceful shutdown work? They are outside the implementation and its [documented deployment evidence](../../HISTORY.md#L50).
- Should oversize frames receive 1009? Today [parser limits](../../main.go#L165) become generic close 1002 through the [handler](../../main.go#L72), while assembled-message overflow explicitly uses 1009.

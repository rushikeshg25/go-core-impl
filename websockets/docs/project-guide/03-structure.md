# Structure

## What lives where

The module has seven tracked files and no runtime subpackages. [main.go](../../main.go) contains the complete HTTP and WebSocket implementation; [main_test.go](../../main_test.go) provides protocol examples and checks, while the script probes upgrade headers. Documentation separates current behavior, intended scope and delivery history.

```text
websockets/
├── go.mod
├── main.go
├── main_test.go
├── ws.sh
├── README.md
├── V1.md
├── HISTORY.md
└── docs/project-guide/   this guide
```

## Module root

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [go.mod](../../go.mod#L1) | Module identity and Go version; no external dependencies. | `module websockets` | Go toolchain |
| [main.go](../../main.go#L1) | HTTP startup, upgrade validation, fragmented echo, masking, lengths, control and close handling. | Exported `WsHandler`, `Health`; local `main`, `token`, `computeWebSocketAcceptKey`, `frame`, `readFrame`, `writeFrame`, `writeClose`, `validClose` | Process runtime, HTTP mux, internal helpers and tests |
| [main_test.go](../../main_test.go#L17) | Raw TCP integration, parser/close/upgrade rejection and fuzz target. | `TestEchoFragmentPingClose`, `TestRejectProtocol`, `TestRejectUpgrade`, `FuzzFrame`; local `masked` | `go test` and Go fuzzing |
| [ws.sh](../../ws.sh#L1) | Curl upgrade probe against localhost:8080. | Shell script | Developer via `bash ws.sh` |
| [README.md](../../README.md#L1) | Current run instructions, supported behavior and exclusions. | Documentation | Maintainers and users |
| [V1.md](../../V1.md#L3) | Scope and acceptance contract. | Documentation | Implementers and reviewers |
| [HISTORY.md](../../HISTORY.md#L10) | Dated commits and delivery verification record. | Documentation | Maintainers checking provenance |

The handler's key internal seam is `readFrame(io.Reader) (frame, error)` followed by message-state handling; output goes through `writeFrame(*bufio.Writer, byte, []byte) error` ([source](../../main.go#L133)). Tests can invoke the handler independently of the fixed production port ([test server](../../main_test.go#L29)).

## Guide folder

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [README.md](README.md) | Entry point, run commands, tour and unresolved questions. | None | New readers |
| [01-architecture.md](01-architecture.md) | Components, contracts and operational limits. | None | Maintainers |
| [02-flow.md](02-flow.md) | Runtime and verification paths. | None | Maintainers |
| [03-structure.md](03-structure.md) | File inventory. | None | Maintainers |
| [04-tech-stack.md](04-tech-stack.md) | Runtime, libraries and tooling. | None | Maintainers |
| [05-decisions.md](05-decisions.md) | Evidence-backed choices and gotchas. | None | Maintainers |

## Excluded

No generated code, vendor directory or lockfile exists in this module. Other modules in the parent checkout are out of scope; its [makefile](../../../makefile#L9) is mentioned only because it contains a WebSocket launch target.

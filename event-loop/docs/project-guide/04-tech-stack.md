# Tech Stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `go 1.23.5` directive; no separate toolchain pin | [go.mod:3](../../go.mod#L3) |
| Module | `github.com/rushikeshg25/event-loop` | [go.mod:1](../../go.mod#L1) |

The package is `main`, so this checkout builds an executable even though its types and methods use exported names ([main.go:1](../../main.go#L1)).

## Frameworks and major libraries

No third-party dependencies or frameworks are declared in [go.mod](../../go.mod#L1). These standard-library packages follow the selected Go toolchain version rather than independent manifest pins.

| Library | Used for | Evidence |
| --- | --- | --- |
| `sync` | Once-only startup, admission mutex, completion wait group | [main.go:20](../../main.go#L20), [main.go:36](../../main.go#L36) |
| `context` | Per-event cooperative cancellation and default Background context | [main.go:45](../../main.go#L45) |
| `fmt` | Demo stdout output | [main.go:138](../../main.go#L138) |
| `testing` | Six behavioral tests | [main_test.go:11](../../main_test.go#L11) |
| `sync/atomic`, `sync`, `time` | Concurrent counters, submitter joins, overlap in tests | [main_test.go:3](../../main_test.go#L3), [main_test.go:13](../../main_test.go#L13) |

Channels, goroutines, and `select` provide the queues and scheduler directly ([main.go:63](../../main.go#L63)). Queue sizes and task limits are literals of five in [construction](../../main.go#L32) and [dispatch](../../main.go#L63), not runtime configuration.

## Data and infrastructure

There is no datastore or external integration: state is channels, closures, and synchronization primitives inside one process ([main.go:10](../../main.go#L10)). There are no environment variables, command-line options, listeners, or service clients in [main](../../main.go#L136).

## Tooling

| Tool | Role | Configured or documented at |
| --- | --- | --- |
| `go run .` | Build and execute the demo | [README.md:3](../../README.md#L3) |
| `go test -race ./...` | Behavioral suite with race detector | [README.md:3](../../README.md#L3), [main_test.go](../../main_test.go#L11) |

The [inventory](03-structure.md#module-root) contains no module build scripts, linter configuration, or CI definition. [HISTORY.md](../../HISTORY.md#delivery-verification) records delivery test success; the guide generation did not rerun it or perform long-running operational validation.

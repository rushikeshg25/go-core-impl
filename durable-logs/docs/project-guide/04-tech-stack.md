# Tech stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `1.23.5` language/toolchain minimum directive, no separate toolchain pin | [go.mod:3](../../go.mod#L3) |
| Protocol Buffers schema | `proto3` syntax | [log.proto:1](../../log.proto#L1) |

## Frameworks and major libraries

| Library | Version | Used for | Evidence |
| --- | --- | --- | --- |
| `google.golang.org/protobuf` | `v1.36.6` | Marshal writes and decode stored records | [go.mod:5](../../go.mod#L5), [writer](../../durablelogs/core.go#L89), [decoder](../../durablelogs/records.go#L83) |
| Go standard library | Supplied by selected Go toolchain | `bufio` buffering, `os` file I/O and sync, `sync.Mutex`, binary framing, timestamps | [core imports](../../durablelogs/core.go#L3), [parser imports](../../durablelogs/records.go#L3) |

The protobuf requirement is marked `indirect` in the manifest despite handwritten direct imports. Treat that annotation as manifest drift, not evidence that protobuf is unused ([go.mod](../../go.mod#L5), [core.go](../../durablelogs/core.go#L9)).

## Data and infrastructure

| Service | Role | Configured at |
| --- | --- | --- |
| Local filesystem | Append-only numbered segment contents, with explicit file sync | [Open arguments and creation](../../durablelogs/core.go#L30), [flush](../../durablelogs/core.go#L110) |
| In-memory buffer and pending slice | Batch writes and expose accepted messages since flush | [state](../../durablelogs/core.go#L18), [pending tracking](../../durablelogs/core.go#L106) |

No external datastore or network dependency appears in the [source inventory](03-structure.md). Directory ownership and file durability limits are described in [Architecture](01-architecture.md#failure-and-scale).

## Tooling

| Tool | Role | Evidence |
| --- | --- | --- |
| `go run .` | Run the demo | [README.md:23](../../README.md#L23), [main.go](../../main.go#L8) |
| `go test -race ./...` | Execute behavioral tests with race detection | [README.md:23](../../README.md#L23), [core_test.go](../../durablelogs/core_test.go#L45) |
| `protoc-gen-go` | Generated bindings report `v1.36.6` | [generated header](../../durablelogs/pb/log.pb.go#L3) |
| `protoc` | Generated bindings report `v5.29.3` | [generated header](../../durablelogs/pb/log.pb.go#L4) |

Generator versions are recorded provenance, not an installed-tool or reproducible-generation pin. No generation command or scoped CI configuration is supplied; see [open questions](README.md#open-questions). The existing [delivery verification](../../HISTORY.md#delivery-verification) reports passing race tests, with no new test execution needed for this documentation-only change.

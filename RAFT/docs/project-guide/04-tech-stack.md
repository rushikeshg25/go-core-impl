# Tech Stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `1.24.6` language/toolchain requirement declared by the module; no separate `toolchain` directive | [go.mod:3](../../go.mod#L3) |
| Go module | `github.com/rushikeshg25/raft`; no dependency requirements | [go.mod:1](../../go.mod#L1) |

## Frameworks and major libraries

There is no application framework or third-party library. The following standard-library packages come with the Go toolchain declared in [go.mod:3](../../go.mod#L3), so they have no independent module version pins.

| Library | Version source | Used for | Used in |
| --- | --- | --- | --- |
| `net`, `net/rpc` | Go toolchain | Raw TCP listeners and synchronous method dispatch/calls; the default RPC codec supplies Go's gob serialization. | [rpc.go:3](../../internal/rpc/rpc.go#L3), [rpc.go:63](../../internal/rpc/rpc.go#L63), [rpc.go:87](../../internal/rpc/rpc.go#L87) |
| `encoding/json`, `os`, `path/filepath` | Go toolchain | Versioned snapshot serialization, directory creation, temporary files, fsync, rename. | [storage.go:3](../../internal/raft/storage.go#L3), [storage.go:74](../../internal/raft/storage.go#L74) |
| `sync`, `time`, `math/rand` | Go toolchain | Mutexes/wait groups, timers/deadlines, randomized election timeout. | [types.go:3](../../internal/raft/types.go#L3), [node.go:3](../../internal/raft/node.go#L3) |
| `flag`, `os/signal`, `syscall` | Go toolchain | CLI configuration and interrupt/termination shutdown. | [main.go:3](../../cmd/raft-demo/main.go#L3), [main.go:53](../../cmd/raft-demo/main.go#L53) |
| `testing` | Go toolchain | Same-package unit/scenario tests and real TCP integration test. | [raft_test.go:3](../../internal/raft/raft_test.go#L3), [rpc_test.go:3](../../internal/rpc/rpc_test.go#L3) |

## Data and infrastructure

| Service or resource | Role | Configured at |
| --- | --- | --- |
| Local filesystem | One exclusively owned JSON state file per persistent node; no external database. | [main.go:18](../../cmd/raft-demo/main.go#L18), [main.go:39](../../cmd/raft-demo/main.go#L39), [storage.go:83](../../internal/raft/storage.go#L83) |
| TCP peer addresses | Ordered static cluster, defaulting to three localhost ports; ID chooses this process's listener. | [main.go:16](../../cmd/raft-demo/main.go#L16), [main.go:48](../../cmd/raft-demo/main.go#L48) |
| RPC timing | 200 ms dial timeout then a 300 ms connection deadline. | [rpc.go:81](../../internal/rpc/rpc.go#L81) |
| Consensus timing | 80 ms heartbeat interval, randomized 600–900 ms election interval, 20 ms ticker. These are source constants/default fields, not CLI flags. | [storage.go:33](../../internal/raft/storage.go#L33), [node.go:10](../../internal/raft/node.go#L10), [node.go:166](../../internal/raft/node.go#L166) |

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| `go run` | Compile/run the node or client mode from this module. | [README.md:7](../../README.md#L7), [main.go:15](../../cmd/raft-demo/main.go#L15) |
| `go test -race ./...` | Package tests with the Go race detector; delivery verification command. | [README.md:35](../../README.md#L35), [V1.md:10](../../V1.md#L10) |
| Injected test transport | Controlled partitions without TCP; directly drives selected protocol transitions. | [raft_test.go:17](../../internal/raft/raft_test.go#L17), [raft_test.go:48](../../internal/raft/raft_test.go#L48) |
| Loopback integration test | Binds OS-assigned ports, starts three real nodes, and observes election/replication. | [rpc_test.go:11](../../internal/rpc/rpc_test.go#L11) |

## Notes

No module-local CI workflow, build script, formatter/linter configuration, or deployment manifest is present in the inspected [file map](03-structure.md#what-lives-where). The version requirement and imports are the available tooling evidence; this guide does not infer parent-repository tooling. [HISTORY.md:74](../../HISTORY.md#L74) records a prior successful race-test run, while [HISTORY.md:82](../../HISTORY.md#L82) explicitly leaves production/long-running operational validation undone. No build or test command was executed while generating this guide.

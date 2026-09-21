# Tech stack

## Languages and runtimes

| Language or runtime | Version | Evidence |
| --- | --- | --- |
| Go | Module declares 1.23.2 | [go.mod:3](../../go.mod#L3) |
| Bash | Unpinned | [ws.sh:1](../../ws.sh#L1) |

## Frameworks and major libraries

All Go libraries are standard-library packages shipped with the toolchain. The [manifest](../../go.mod#L1) has no third-party requirements.

| Library | Version | Used for | Used in |
| --- | --- | --- | --- |
| `net/http`, `net/url` | Go toolchain | Routing, HTTP server, hijack, origin parsing | [main.go:34](../../main.go#L34), [main.go:216](../../main.go#L216) |
| `crypto/sha1`, `encoding/base64` | Go toolchain | RFC upgrade accept calculation and key decoding | [main.go:30](../../main.go#L30) |
| `bufio`, `io`, `encoding/binary` | Go toolchain | Buffered wire I/O, complete reads, big-endian frame lengths and close codes | [main.go:133](../../main.go#L133), [main.go:184](../../main.go#L184) |
| `unicode/utf8`, `time` | Go toolchain | Text/close validation and socket deadlines | [main.go:69](../../main.go#L69), [main.go:114](../../main.go#L114), [main.go:205](../../main.go#L205) |
| `testing`, `net/http/httptest`, `net` | Go toolchain | Unit/fuzz testing and local HTTP/TCP integration | [main_test.go:3](../../main_test.go#L3) |

## Data and infrastructure

| Resource | Role | Configured at |
| --- | --- | --- |
| TCP listener on `:8080` | Plain HTTP and upgraded WebSocket connections | [main.go:220](../../main.go#L220) |
| Connection-local memory | Frame payload and fragmented-message assembly, each bounded by 1 MiB checks | [main.go:20](../../main.go#L20), [main.go:108](../../main.go#L108), [main.go:165](../../main.go#L165) |

No datastore or external API is used by the [runtime source](../../main.go#L3). Authentication, TLS, compression, extensions and subprotocol negotiation are explicitly outside the [README's v1 behavior](../../README.md#L7).

## Tooling

| Tool | Role | Evidence |
| --- | --- | --- |
| `go run .` | Compile and launch the module | [README.md:3](../../README.md#L3) |
| `go test -race ./...` | Execute tests with race detection | [README.md:11](../../README.md#L11), [main_test.go:28](../../main_test.go#L28) |
| Go native fuzzing | Feed arbitrary byte streams to the parser | [README.md:11](../../README.md#L11), [main_test.go:86](../../main_test.go#L86) |
| curl | Send upgrade request headers, version unpinned | [ws.sh:6](../../ws.sh#L6) |
| Parent Makefile | Has a launch target whose `make go run .` does not match the documented Go command | [makefile:9](../../../makefile#L9) |

The tracked module inventory contains no CI, container, deployment or lint configuration ([inventory](03-structure.md#module-root)). [HISTORY.md](../../HISTORY.md#L42) records a passing race-test run; it also explicitly records the absence of production or long-running operational validation.

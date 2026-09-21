# RFC 6455 echo server history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2024-12-14T07:46:35+05:30: init: Ws

- **What happened:** The repository records `init: Ws`.
- **Evidence:** [commit c140b2763f](https://github.com/rushikeshg25/go-core-impl/commit/c140b2763fbe95a490ee2fe1697727772f880261).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:21:37+05:30: docs: define websockets v1 contract

- **What happened:** The repository records `docs: define websockets v1 contract`.
- **Evidence:** [commit 7e620439e4](https://github.com/rushikeshg25/go-core-impl/commit/7e620439e472be888c2a4dde5d2fdbaf350519a9).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:37:41+05:30: feat(websockets): implement bounded RFC6455 echo and control frames

- **What happened:** The repository records `feat(websockets): implement bounded RFC6455 echo and control frames`.
- **Evidence:** [commit 1752eb2548](https://github.com/rushikeshg25/go-core-impl/commit/1752eb2548b8ecf0e22b27d14e53c0b52d808db6).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:37:41+05:30: test(websockets): verify raw upgrade fragmentation control frames and malformed input

- **What happened:** The repository records `test(websockets): verify raw upgrade fragmentation control frames and malformed input`.
- **Evidence:** [commit 57ca7453bd](https://github.com/rushikeshg25/go-core-impl/commit/57ca7453bdfe0447d4d5643ef45f27555965413d).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:39:40+05:30: docs: document websockets v1 usage and limitations

- **What happened:** The repository records `docs: document websockets v1 usage and limitations`.
- **Evidence:** [commit 1edf6755a2](https://github.com/rushikeshg25/go-core-impl/commit/1edf6755a2bd4658e6746583a63f0f57232c95d0).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed `go test -race ./...` with real HTTP/TCP upgrade, fragmented echo, ping/pong and close exchange in [main_test.go](main_test.go).

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.

## 2026-09-21: Maintainer project guide

Added the six-file [project guide](docs/project-guide/README.md), tracing architecture, runtime flows, source structure, dependencies and decisions against the v1 code. Relative paths, source/heading anchors and Mermaid syntax were checked. The guide distinguishes observed behavior from inferred rationale and records remaining limitations.

# Upload-to-HLS service history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2025-02-08T20:19:48+05:30: init:hls

- **What happened:** The repository records `init:hls`.
- **Evidence:** [commit 4ff710e67a](https://github.com/rushikeshg25/go-core-impl/commit/4ff710e67a2c764179078d445482e6f9586fea52).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:21:37+05:30: docs: define hls v1 contract

- **What happened:** The repository records `docs: define hls v1 contract`.
- **Evidence:** [commit 69552e8ddd](https://github.com/rushikeshg25/go-core-impl/commit/69552e8ddd802319b5fc197ef0d2574211732a1a).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:46:28+05:30: feat(hls): implement isolated uploads processing status and bounded transcoding

- **What happened:** The repository records `feat(hls): implement isolated uploads processing status and bounded transcoding`.
- **Evidence:** [commit beb7b49047](https://github.com/rushikeshg25/go-core-impl/commit/beb7b490476d7b58e8dceffab0b006ffcae01c00).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:46:28+05:30: feat(hls): connect upload progress processing states and playlist playback

- **What happened:** The repository records `feat(hls): connect upload progress processing states and playlist playback`.
- **Evidence:** [commit 1011a0d789](https://github.com/rushikeshg25/go-core-impl/commit/1011a0d7897af10c3c886cca660e516748443375).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:46:28+05:30: test(hls): verify upload status restart and real FFmpeg playlists

- **What happened:** The repository records `test(hls): verify upload status restart and real FFmpeg playlists`.
- **Evidence:** [commit 344520fa82](https://github.com/rushikeshg25/go-core-impl/commit/344520fa827a9456082b7e9f1866eec1af92771f).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:59:21+05:30: fix(hls): drain accepted uploads on shutdown and lock client dependencies

- **What happened:** The repository records `fix(hls): drain accepted uploads on shutdown and lock client dependencies`.
- **Evidence:** [commit b139f40d91](https://github.com/rushikeshg25/go-core-impl/commit/b139f40d91ea528e8862d464b6718d413315ed18).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T14:04:58+05:30: docs: document hls v1 usage and limitations

- **What happened:** The repository records `docs: document hls v1 usage and limitations`.
- **Evidence:** [commit e3c69d7c2f](https://github.com/rushikeshg25/go-core-impl/commit/e3c69d7c2ffdf3a9337c302d73bd66d637594687).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed server `go test -race ./...`, including real FFmpeg generation in [main_test.go](server/main_test.go). Client `npm run build` and `npm run lint` passed using the committed package lock.

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.

## 2026-09-21: Maintainer project guide

Added the six-file [project guide](docs/project-guide/README.md), tracing architecture, runtime flows, source structure, dependencies and decisions against the v1 code. Relative paths, source/heading anchors and Mermaid syntax were checked. The guide distinguishes observed behavior from inferred rationale and records remaining limitations.

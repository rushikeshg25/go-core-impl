# V1 delivery integration

This integration links fourteen independent v1 project PRs. Nine are separate repositories pinned as submodules; five are directories with their own PRs in this repository. No PR is merged by this delivery.

## Review and merge order

Review each project contract, README, HISTORY.md and docs/project-guide/. Each project has its own branch and 5–10 commits. Merge the submodule PRs before this pin PR, preserving their commit heads with a merge commit. If a project is squash-merged or rebased, refresh its pin to the resulting commit before merging this integration. The five directory PRs can be reviewed independently; merge them before this integration so the updated catalog matches the source. This integration does not include their source changes.

## Project PRs

| Project | PR | Location |
| --- | --- | --- |
| append-only-storage | [V1 PR](https://github.com/rushikeshg25/append-only-storage/pull/1) | submodule |
| adaptive-bitrate-streaming | [V1 PR](https://github.com/rushikeshg25/adaptive-bitrate-streaming/pull/1) | submodule |
| bp-tree | [V1 PR](https://github.com/rushikeshg25/bp-tree/pull/1) | submodule |
| loadbalancer | [V1 PR](https://github.com/rushikeshg25/loadbalancer/pull/2) | submodule |
| p2p-file-sharing | [V1 PR](https://github.com/rushikeshg25/p2p-file-sharing/pull/1) | submodule |
| task-scheduler | [V1 PR](https://github.com/rushikeshg25/task-scheduler/pull/1) | submodule |
| token-bucket | [V1 PR](https://github.com/rushikeshg25/token-bucket/pull/1) | submodule |
| tricolor-gc | [V1 PR](https://github.com/rushikeshg25/tricolor-gc/pull/1) | submodule |
| wal-go | [V1 PR](https://github.com/rushikeshg25/wal-go/pull/1) | submodule |
| event-loop | [V1 PR](https://github.com/rushikeshg25/go-core-impl/pull/1) | directory |
| websockets | [V1 PR](https://github.com/rushikeshg25/go-core-impl/pull/2) | directory |
| hls | [V1 PR](https://github.com/rushikeshg25/go-core-impl/pull/4) | directory |
| RAFT | [V1 PR](https://github.com/rushikeshg25/go-core-impl/pull/5) | directory |
| durable-logs | [V1 PR](https://github.com/rushikeshg25/go-core-impl/pull/3) | directory |

## Verification

All fourteen project modules passed `go test -race ./...` for their v1 code. The video clients passed their production builds and lint checks. FFmpeg was installed and exercised by the video tests. Raft was tested both through deterministic partition/restart scenarios and a real localhost three-node RPC cluster. Test scope and v1 limitations are recorded in the individual histories and READMEs.

The pin commits record exact project branch heads. Project guides are included in those heads and use relative links and Mermaid diagrams. WAL v1 deliberately rejects legacy unversioned files; use a new directory or an explicit migration process.

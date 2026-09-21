# V1 integration history

> Scope: parent repository integration of fourteen independent project deliveries.
> Last updated: 2026-09-21. Git timestamps retain their author timezone offsets.

## At a glance

Nine submodule projects and five directory projects now have individual v1 review branches with 5–10 commits, executable acceptance checks, evidence-backed histories, and six-file project guides. This branch pins the nine submodules and records the full review set. The directory implementations remain in their separate PRs.

## Timeline

### 2026-09-18: Storage series baseline

The parent recorded storage-compaction/tombstone updates and checksum-corruption-detection, with a refreshed catalog. Evidence: [660bd23](https://github.com/rushikeshg25/go-core-impl/commit/660bd231d7d459aab1fccee5cfbfd37308aa5712) and [4a58bec](https://github.com/rushikeshg25/go-core-impl/commit/4a58bece55aaa47598f912881f3c6b62a1ac0ae6). Confirmed by Git history; these commits establish the starting parent state, not the fourteen v1 implementations.

### 2026-09-21: Independent v1 delivery and source guides

Each project received a bounded v1 contract, implementation and behavioral checks, documentation, and history. Fresh project-scoped guide agents read the completed source; their findings were checked and guide links/anchors and Mermaid syntax validated. All fourteen Go modules passed the race suite, both video clients passed build/lint, installed FFmpeg ran in video tests, and Raft exercised a real three-node TCP cluster. Evidence and limitations live in each project's HISTORY.md and guide; [the delivery table](docs/v1-delivery.md) links every PR.

### 2026-09-21: Final submodule pins

The parent pins the nine final project branch heads including guides, grouped into storage, networking, and data structure/concurrency commits. [The manifest](docs/v1-delivery.json) records immutable SHAs and commit counts. Remote branch heads were compared before indexing the gitlinks. Merge submodule PRs first while preserving these commits; after squash/rebase, refresh pins to the resulting heads before integrating.

## Turning points

The v1 contracts bounded each project's promises and made failure and lifecycle behavior reviewable. Source-grounded guides exposed additional adaptive streaming publication/client retry gaps, a Raft failure-snapshot gap, and missing event-loop cancellation; their fixes and rerun checks are recorded in those projects' histories.

## Open questions

No PR was merged and no production deployment was performed by this delivery. Review acceptance and operational limits per project. WAL v1 rejects legacy unversioned files and includes no migration tool.

## Sources

- [Delivery PRs and merge order](docs/v1-delivery.md)
- [Final project commits](docs/v1-delivery.json)
- Each linked project's V1.md, HISTORY.md, tests, and docs/project-guide/.

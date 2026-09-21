# Structure

## What lives where

The module contains 14 tracked files at the guide's source revision, principally Go source plus one protobuf schema. The root supplies a demo and module contract; `durablelogs/` contains the handwritten library and tests, while `durablelogs/pb/` contains generated wire types. This guide is scoped to this module, not sibling experiments in the enclosing repository.

```text
main.go                 demo executable
go.mod                  module and dependencies
log.proto               stored record schema
README.md / V1.md        API guide and acceptance contract
HISTORY.md              delivery evidence
durablelogs/
  core.go               logger state and lifecycle
  records.go            discovery and decoding
  util.go               legacy panic helpers
  core_test.go          behavioral tests
  pb/                   generated protobuf bindings
docs/project-guide/     this guide
```

## Module root

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [main.go](../../main.go#L8) | Open local demo logger, append ten messages, close with error checks. | `main` | `go run .` |
| [log.proto](../../log.proto#L1) | Proto3 message and timestamp fields. | `Log` schema | Generated Go binding, external generation tooling |
| [go.mod](../../go.mod#L1) | Module identity, Go directive, protobuf dependency. | Module `durablelogs` | Go tooling |
| [README.md](../../README.md#L1) | Public usage and storage/failure limits. | None | Maintainers and callers |
| [V1.md](../../V1.md#L3) | v1 scope and acceptance criteria. | None | Delivery and review |
| [HISTORY.md](../../HISTORY.md#L54) | Historical delivery and test evidence. | None | Maintainers |

## Durablelogs

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [core.go](../../durablelogs/core.go#L15) | Own writer state, mutex, errors, pending snapshots, rotation and replay. | `DurableLogger`, `Open`, `NewDLServer`, `ErrClosed`, `ErrPoisoned`; methods `Log`, `Flush`, `NewFile`, `GetBufferedLogs`, `GetCurrentFile`, `ReadAll`, `Close` | Demo and library consumers |
| [records.go](../../durablelogs/records.go#L17) | Validate and sort segment names, frame size bounds, decode, optional tail repair. | None; internal `segments`, `segmentPath`, `readSegment`, `maxRecord` | `Open`, `ReadAll`, rotation, tests |
| [util.go](../../durablelogs/util.go#L9) | Legacy protobuf helpers that panic on errors. | `MustMarshal`, `MustUnmarshal` | No handwritten call sites in this module |
| [core_test.go](../../durablelogs/core_test.go#L10) | Four tests covering rotation/reopen/snapshot ownership, concurrent appends, tail recovery/closed state, and poisoned flush. | `TestRotateReopenAndBuffer`, `TestConcurrentLog`, `TestTailRecoveryAndClosed`, `TestPoisonOnFlush` | `go test` |

## Generated bindings

Generated implementation details are excluded from the handwritten source inventory. These two entries locate the runtime dependency and flag the unexplained legacy artifact.

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [log.pb.go](../../durablelogs/pb/log.pb.go#L24) | Generated `log.proto` representation and descriptors. | `Log`, `File_log_proto` | Logger, decoder, marshal helpers |
| [message.pb.go](../../durablelogs/pb/message.pb.go#L24) | Generated legacy message representation, source schema absent. | `Mesage`, `File_message_proto` | No handwritten call sites in this module |

## Excluded

The dependency checksum file and ignore boilerplate do not add design information. Generated descriptor internals are omitted above. No scoped CI, build script, container configuration, vendor tree, or deployment manifest is present.

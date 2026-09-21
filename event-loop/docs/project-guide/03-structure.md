# Structure

## What lives where

The scoped module has seven tracked files at the documented commit: two Go files, one manifest, three Markdown files, and a generic ignore file. Production logic and the demo share [main.go](../../main.go); tests share its `main` package in [main_test.go](../../main_test.go). Other repository modules are outside this guide.

```text
event-loop/
  main.go             # event types, admission, dispatcher, lifecycle, demo
  main_test.go        # six behavioral tests
  go.mod              # module identity and Go version
  README.md           # supported usage and limitations
  V1.md               # intended contract and acceptance criteria
  HISTORY.md          # delivery history and verification record
  .gitignore          # generic Go exclusions
  docs/project-guide/ # this six-file guide
```

## Module root

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [main.go](../../main.go#L1) | Entire runtime and runnable example | `Event`, `EventLoop`, `NewEventLoop`, `Start`, `AddEvent`, `StopEventLoop`, `AddCallback` | Go executable entry point and tests; dispatcher calls private `runTask` and `invoke` |
| [main_test.go](../../main_test.go#L11) | Drain and bound check, concurrent admission/stop, synchronous ordering, cancellation, invalid admission | `TestDrainAndBound`, `TestConcurrentSubmitStop`, `TestSyncCallbackOrder`, `TestCancellationStillCompletesCallbacks`, `TestRunningTaskCooperatesWithCancellation`, `TestInvalidAndStoppedCallbackAdmission` | `go test` |
| [go.mod](../../go.mod#L1) | Module identity and minimum Go language/toolchain requirement | None | Go tooling |
| [README.md](../../README.md#L1) | Commands, public behavior, unsupported operations | None | Maintainers |
| [V1.md](../../V1.md#L1) | Desired v1 behavior and acceptance | None | Maintainers; compare [open questions](README.md#open-questions) |
| [HISTORY.md](../../HISTORY.md#L1) | Commit chronology and delivery verification | None | Maintainers |

The test suite uses atomics for concurrent counters, a wait group for concurrent submitters, and a small sleep to let tasks overlap ([main_test.go:13](../../main_test.go#L13)). The submit/stop test exercises concurrency but does not assert accepted-work counts ([main_test.go:38](../../main_test.go#L38)); race checking adds useful runtime validation. Tests are not proof against all deadlocks.

## Guide folder

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [README.md](README.md) | Tour, run commands, unresolved gaps | None | New maintainers |
| [01-architecture.md](01-architecture.md) | Components, state, boundaries | None | Guide index |
| [02-flow.md](02-flow.md) | Runtime flow diagrams and source steps | None | Guide index |
| [03-structure.md](03-structure.md) | File inventory | None | Other guide pages |
| [04-tech-stack.md](04-tech-stack.md) | Go and standard-library tooling | None | Guide index |
| [05-decisions.md](05-decisions.md) | Choices, rationale confidence, gotchas | None | Guide index |

## Excluded

Only the boilerplate [.gitignore](../../.gitignore) is excluded from the responsibility tables. There are no generated source files, vendored dependencies, lockfiles, module-specific CI, containers, or deployment definitions in this scoped module.

# Decisions

These choices are visible in code. Rationale is marked inferred unless explicitly documented.

## Serialize callbacks through one coordinator

- **What:** Synchronous tasks and all callbacks execute in `run`; asynchronous task bodies run in separate goroutines.
- **Evidence:** [main.go:76](../../main.go#L76), [README.md:5](../../README.md#L5).
- **Why, apparently:** Inferred: keep callbacks mutually serialized while allowing task work to overlap.
- **Tradeoff:** A slow callback pauses queue processing and completion handling. Serialization does not prevent races with asynchronous task bodies.
- **Confidence:** Execution model confirmed by code and README; motivation inferred.

## Bound dispatch with a counter rather than reusable workers

- **What:** At five outstanding asynchronous events, the coordinator disables its events receive case. Every admitted asynchronous task still gets a newly launched goroutine.
- **Evidence:** [main.go:63](../../main.go#L63), [main.go:71](../../main.go#L71), [main.go:79](../../main.go#L79).
- **Why, apparently:** Inferred: small, explicit concurrency bound without a worker lifecycle.
- **Tradeoff:** Synchronous events behind this queue also wait; queued completions occupy active slots until consumed. Capacity is hard-coded.
- **Confidence:** Mechanism confirmed; rationale inferred. The [Event comment](../../main.go#L9) calls this a bounded pool, but there are no persistent worker goroutines.

## Keep admission and stopping under one mutex

- **What:** Submission checks stopped state and sends while holding `mu`; stop takes the same lock before closing its signal channel.
- **Evidence:** [main.go:103](../../main.go#L103), [main.go:115](../../main.go#L115), [main.go:128](../../main.go#L128).
- **Why, apparently:** Inferred: ensure a successful submission is queued before shutdown can stop admission.
- **Tradeoff:** Full queues hold the mutex while waiting, delaying other producers and shutdown. Reentrant submission from coordinator code can deadlock.
- **Confidence:** Ordering confirmed by code; intent inferred and blocking restriction documented in [README.md:9](../../README.md#L9).

## Use a single-use draining lifecycle

- **What:** `sync.Once` starts one coordinator; stop closes a signal once and every caller waits on `done`. The loop exits only when queues and outstanding tasks are empty.
- **Evidence:** [main.go:36](../../main.go#L36), [main.go:68](../../main.go#L68), [main.go:113](../../main.go#L113).
- **Why, apparently:** Graceful drain is explicit in [README.md:7](../../README.md#L7); inferred benefit is simple repeated-stop behavior.
- **Tradeoff:** There is no restart or loop shutdown timeout. Stop does not cancel per-event contexts; callers own cancellation and task functions must cooperate ([runTask](../../main.go#L45)).
- **Confidence:** Drain confirmed by implementation and README.

## Make cancellation cooperative and preserve completion callbacks

- **What:** An optional event context can skip a not-yet-started task; `TaskContext` receives that context and takes precedence over legacy `Task`. A canceled event still reaches its callback.
- **Evidence:** [main.go:45](../../main.go#L45), [dispatch](../../main.go#L76), [README.md:13](../../README.md#L13).
- **Why:** The [history follow-up](../../HISTORY.md#2026-09-21-cooperative-event-cancellation) records this as closing the cancellation gap found against the v1 contract.
- **Tradeoff:** Context-aware tasks choose how to stop. A legacy task already running remains noninterruptible; callbacks receive no cancellation status argument.
- **Confidence:** Behavior and delivery rationale documented; the [cancellation tests](../../main_test.go#L63) exercise skipped work and running cooperation.

## Gotchas

- Construct with [NewEventLoop](../../main.go#L31); a zero-value loop has nil channels and is not usable for submission or stop.
- Call submission and stop from outside tasks/callbacks. A callback waiting for shutdown prevents its own coordinator from finishing; queue saturation can similarly deadlock nested submission ([README.md:9](../../README.md#L9)).
- [AddCallback](../../main.go#L123) executes an event's task too, always on the coordinator, regardless of `Async` ([dispatch](../../main.go#L87)).
- [Start](../../main.go#L36) returns a wait group for the loop's lifetime, not for the currently queued batch. Call stop before waiting for final completion; `done` closes just before deferred `wg.Done` ([main.go:61](../../main.go#L61)).
- Nil functions are safe, but panic recovery is absent ([invoke](../../main.go#L40)). Callback guarantees assume normal task return, a live process, and supported submission APIs.
- Public channels are legacy exposure, not an invitation to close or send directly ([README.md:9](../../README.md#L9)). Doing so bypasses stopped checks and can invalidate draining behavior.

## Conventions

Lifecycle tests live beside the implementation in the same package and exercise exported operations ([main_test.go:1](../../main_test.go#L1)). Preserve task-before-callback order and check concurrency changes with the documented race-enabled test command ([README.md:3](../../README.md#L3)); see [coverage gaps](README.md#open-questions) before treating the current suite as exhaustive.

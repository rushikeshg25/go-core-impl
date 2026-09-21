# Event loop v1

Run `go run .` for a task/callback demo and `go test -race ./...` for lifecycle tests.

`NewEventLoop` creates a loop. `Start` is idempotent and returns a wait group immediately. `AddEvent` starts the loop if needed and returns whether the event was accepted. Set `Async: true` for asynchronous work; at most five such tasks run at once. Synchronous tasks and callbacks run on the coordinator, in task-before-callback order. Async completion order is unspecified.

`StopEventLoop` stops admission, drains accepted events and callbacks, then waits for completion. Repeated and concurrent stops are safe. A task must return for shutdown to complete; Go functions cannot be forcibly cancelled.

Submit work and stop the loop from outside tasks/callbacks: these operations can block and calling them inside the coordinator can deadlock. Use `AddEvent` and `AddCallback`; direct writes or closure of the legacy exported channels are unsupported. Task panics retain normal Go process-panic behavior.

## Cooperative cancellation

Set `Event.Context` and `TaskContext: func(context.Context)` for cancellable work. `TaskContext` takes precedence over `Task`; a nil context defaults to `context.Background()`. If the context is already canceled when execution begins, the task is skipped and its callback still runs once. Cancellation during execution is cooperative: `TaskContext` must observe `ctx.Done()` and return before the callback and draining shutdown can complete. Legacy `Task: func()` cannot be interrupted once running. Cancellation does not retract admission or suppress callbacks.

# Decisions

## Own the protocol codec

- **What:** Implement upgrade and frame processing directly using Go's standard library.
- **Evidence:** [go.mod](../../go.mod#L1) has no external dependencies; [main.go](../../main.go#L30) computes the accept key, reads masking/length fields and writes frames.
- **Why, apparently:** Making protocol mechanics visible fits the [v1 contract](../../V1.md#L5). This educational rationale is inferred, not stated by an author comment.
- **Tradeoff:** A compact implementation is easy to inspect, but protocol compliance must be established by its own tests. The current [test suite](../../main_test.go#L28) covers only selected paths.
- **Confidence:** Implementation confirmed; rationale inferred.

## Keep the hijacker's buffered reader

- **What:** Read frames from `rw.Reader` returned by `Hijack`, instead of constructing a new reader around the raw socket.
- **Evidence:** [main.go:57](../../main.go#L57), [main.go:70](../../main.go#L70); the [README](../../README.md#L5) explicitly promises preservation of buffered upgrade bytes.
- **Why:** A client may send frame bytes alongside the upgrade request, leaving them already buffered by HTTP processing.
- **Tradeoff:** The connection loop now owns flushing, deadlines and closure; ordinary HTTP response handling ends at hijack.
- **Confidence:** Confirmed by code and README.

## Assemble messages before echo

- **What:** Preserve one initial opcode, accumulate continuation payloads and return a single final frame after FIN. Control frames bypass assembly.
- **Evidence:** [main.go:86](../../main.go#L86), [main.go:93](../../main.go#L93), [main.go:113](../../main.go#L113), with a ping-between-fragments [test](../../main_test.go#L46).
- **Why, apparently:** Complete-message validation permits UTF-8 code points to span fragments; this rationale is inferred from validation after assembly.
- **Tradeoff:** Fragment boundaries are lost, memory grows with message size, and no echo is produced until completion. The aggregate 1 MiB cap bounds message size.
- **Confidence:** Behavior confirmed; rationale inferred.

## Bound each connection with sizes and deadlines

- **What:** Limit data frames and messages to 1 MiB and controls to 125 bytes, with 60-second read and 10-second write deadlines.
- **Evidence:** [main.go:20](../../main.go#L20), [main.go:69](../../main.go#L69), [main.go:108](../../main.go#L108), [main.go:165](../../main.go#L165).
- **Why, apparently:** Bound malformed or slow-client resource use without a shared session manager; inferred from the checks.
- **Tradeoff:** These are per-connection constraints, not total memory/connection limits. The read deadline resets per frame, so timely fragments or controls can keep an unfinished message alive.
- **Confidence:** Constraints confirmed; rationale inferred.

## Allow clients without Origin

- **What:** Validate HTTP(S) origin and matching host only when the header is present.
- **Evidence:** [main.go:45](../../main.go#L45) and the explicit [README contract](../../README.md#L5).
- **Why, apparently:** Support non-browser clients as well as same-host browsers; rationale inferred.
- **Tradeoff:** Origin checking is not authentication, and it does not prevent a non-browser client from omitting or supplying that header.
- **Confidence:** Policy confirmed; rationale inferred.

## Gotchas

- **Close codes differ by error location:** a frame over 1 MiB becomes generic 1002, whereas assembled overflow becomes 1009. A read timeout also goes through 1002. The README's general rejection description omits these distinctions ([parser](../../main.go#L165), [handler](../../main.go#L72), [assembly](../../main.go#L108)).
- **Deadlines are frame based:** the README calls them idle-read timeouts, but one read deadline covers the complete `readFrame` call and is not refreshed for every byte ([main.go:69](../../main.go#L69)).
- **Close is best effort:** generated closes discard write errors; valid incoming closes are echoed and the socket is immediately closed. There is no server-initiated close wait or graceful connection drain ([main.go:79](../../main.go#L79), [main.go:202](../../main.go#L202), [startup](../../main.go#L216)).
- **The script is an upgrade probe:** [ws.sh](../../ws.sh#L6) sends headers but does not construct frames. Its public sample key is an RFC handshake example, not authentication.
- **Tests do not imply full conformance:** the [integration helper](../../main_test.go#L17) supports short payloads only, and [fuzzing](../../main_test.go#L86) targets parser robustness rather than semantic assertions. Missing cases are tracked in [Open questions](README.md#open-questions).
- **Run from the module directory:** the parent [makefile](../../../makefile#L9) calls `make go run .`; there is no module Makefile defining those targets. The documented command is [go run .](../../README.md#L3).

## Conventions

The runtime and tests share `package main`, so tests directly exercise unexported codec functions ([runtime](../../main.go#L1), [tests](../../main_test.go#L1)). Protocol errors return early, while connection cleanup is centralized in a deferred close ([handler](../../main.go#L61)). Match the existing distinction between byte-level frame checks in `readFrame` and cross-frame message checks in `WsHandler` when adding behavior ([parser](../../main.go#L133), [assembly](../../main.go#L93)).

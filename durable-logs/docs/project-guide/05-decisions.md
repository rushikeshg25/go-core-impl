# Decisions

Implemented behavior below is evidenced directly; rationale is marked inferred when the source does not state it.

## Keep framing simple and compatible

- **What:** Prefix each protobuf `Log` with a four-byte little-endian body length and retain message/timestamp fields.
- **Evidence:** [Writer](../../durablelogs/core.go#L98), [decoder](../../durablelogs/records.go#L63), and [schema](../../log.proto#L7); the [README](../../README.md#L17) explicitly says the existing format is preserved.
- **Why, apparently:** Compatibility is documented. Straightforward sequential decoding is an inferred benefit.
- **Tradeoff:** There is no checksum, magic header, or explicit format-version field, so arbitrary byte corruption may still decode successfully.
- **Confidence:** Format and compatibility confirmed; simplicity rationale inferred.

## Separate accepted messages from durable flushes

- **What:** Track pending message strings independently of the bufio buffer and clear them only after flush plus file sync.
- **Evidence:** [Pending append](../../durablelogs/core.go#L106), [flush](../../durablelogs/core.go#L110), [owned copy](../../durablelogs/core.go#L155), [snapshot test](../../durablelogs/core_test.go#L21).
- **Why, apparently:** The [README](../../README.md#L15) explicitly defines retrieval as messages accepted since the last successful flush, even after automatic bufio draining.
- **Tradeoff:** Callers get a meaningful pending snapshot, but strings consume additional memory and `Log` success alone does not promise durability.
- **Confidence:** Confirmed by code, tests, and documentation.

## Recover only the latest incomplete tail

- **What:** Validate every segment on open and allow truncation only for an incomplete record in the highest-numbered segment.
- **Evidence:** [Repair flag selection](../../durablelogs/core.go#L42), [bounded parsing and truncation](../../durablelogs/records.go#L59).
- **Why, apparently:** Inferred: an interrupted append is plausible at the active tail; silently discarding older or malformed complete records could conceal corruption.
- **Tradeoff:** Restarts recover partial final appends, but still fail for malformed protobuf, oversized lengths, earlier truncation, and ordinary I/O errors. Startup must scan historical data.
- **Confidence:** Behavior confirmed; rationale inferred.

## Serialize the whole lifecycle and stop after I/O failure

- **What:** Public instance methods lock one mutex. Failed write, flush, or rotation records poison; guarded operations reject later calls, while close attempts cleanup.
- **Evidence:** [State checks](../../durablelogs/core.go#L65), [Log](../../durablelogs/core.go#L80), [ReadAll](../../durablelogs/core.go#L165), [Close](../../durablelogs/core.go#L188), [poison test](../../durablelogs/core_test.go#L87).
- **Why, apparently:** Inferred: serialize file/count transitions and avoid continuing appends after an uncertain write boundary.
- **Tradeoff:** Simple concurrency semantics block writers during full replay and sync. Separate instances sharing one directory are still unsafe.
- **Confidence:** Behavior confirmed; rationale inferred.

## Rotate by record count and preserve legacy entry points

- **What:** Rotate before the next append once the current count reaches the limit; retain explicit `NewFile`, panic-on-open `NewDLServer`, and panic marshal helpers.
- **Evidence:** [Count check](../../durablelogs/core.go#L93), [exclusive next file](../../durablelogs/core.go#L133), [legacy constructor comment](../../durablelogs/core.go#L57), [helpers](../../durablelogs/util.go#L9).
- **Why, apparently:** Legacy constructor compatibility is explicitly stated; count-based sizing keeps the demo predictable (inferred).
- **Tradeoff:** Segment byte sizes vary with messages, explicit rotation can create empty segments, and panic helpers do not follow the main error-returning API.
- **Confidence:** Mechanics and constructor status confirmed; count rationale inferred.

## Gotchas

- `ReadAll` changes persistence state: it flushes and clears pending messages before reading ([core.go:171](../../durablelogs/core.go#L171)).
- Rotation happens on the next write, so a file exactly at capacity remains current until another append; reopening preserves that behavior ([core.go:48](../../durablelogs/core.go#L48), [core.go:93](../../durablelogs/core.go#L93)).
- Names such as `dl-01`, `dl-bad`, and a directory named `dl-0` fail discovery; unrelated names are ignored, and gaps between IDs are accepted ([records.go:25](../../durablelogs/records.go#L25)).
- The writer's message limit is `64 MiB - 1024` bytes, while the decoder bounds the encoded payload at 64 MiB ([core.go:86](../../durablelogs/core.go#L86), [records.go:17](../../durablelogs/records.go#L17)).
- `GetBufferedLogs` and `GetCurrentFile` do not reject closed or poisoned instances. Repeated `Close` returns nil, so preserve the first close error ([getters](../../durablelogs/core.go#L155), [Close](../../durablelogs/core.go#L188)).
- File sync does not sync directory entries. A second instance is not blocked from opening the same directory ([README](../../README.md#L19), [Open](../../durablelogs/core.go#L30)).
- The generated `Mesage` artifact and schema import-path mismatch need clarification before establishing a regeneration workflow ([open questions](README.md#open-questions)).

## Conventions

Use the error-returning `Open` and check lifecycle errors, following the [demo](../../main.go#L8). Keep storage helpers private and expose caller operations through `DurableLogger`, as in [records.go](../../durablelogs/records.go#L19) and [core.go](../../durablelogs/core.go#L80). Tests use package `durablelogs` and temporary directories so they can exercise internal failure state without persistent fixtures ([core_test.go](../../durablelogs/core_test.go#L1)).

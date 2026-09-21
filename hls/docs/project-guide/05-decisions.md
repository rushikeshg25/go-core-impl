# Decisions

The evidence below distinguishes documented scope from rationale inferred from implementation. See [V1.md](../../V1.md#L5) for the delivery contract and [HISTORY.md](../../HISTORY.md#L24) for the recorded implementation sequence.

## Admit one upload through processing

- **What:** A single channel slot is acquired before parsing an upload and released after its worker finishes; saturation returns 503 immediately.
- **Evidence:** [capacity](../../server/main.go#L49), [acquisition](../../server/main.go#L166), [worker release](../../server/main.go#L232).
- **Why:** Bounded processing and retry-on-busy are explicit in [README.md](../../README.md#L20). Inferred: including receipt/parsing in the slot prevents concurrent uploads from accumulating disk work before encoder admission.
- **Tradeoff:** Resource use is simpler to bound, but a slow upload blocks admission for every other user and there is no queued retry.
- **Confidence:** Policy documented; reason for acquisition timing inferred.

## Use opaque storage IDs and constrained file serving

- **What:** Generate a 16-byte random ID, use it as the storage directory, retain only a basename as the display name, and serve only completed playlists/segments.
- **Evidence:** [ID and directory](../../server/main.go#L194), [display name](../../server/main.go#L222), [media gate](../../server/main.go#L127).
- **Why:** Upload isolation is an explicit [contract](../../V1.md#L5); the [README](../../README.md#L20) distinguishes display filenames from storage paths.
- **Tradeoff:** Simple separation and restrictive routing do not provide identity or authorization. Anyone who can reach the collection route can list all video IDs ([listing](../../server/main.go#L106)).
- **Confidence:** Documented and implemented.

## Store metadata beside generated media

- **What:** Persist one JSON record per directory by temp-file rename and maintain a map for request reads. Restore the map on startup and convert interrupted processing to failure.
- **Evidence:** [restore](../../server/main.go#L44), [save](../../server/main.go#L76).
- **Why, apparently:** Inferred: keeping status with media supplies restart visibility without a database and fits the documented single-process scope ([README](../../README.md#L22)).
- **Tradeoff:** Rename avoids exposing partially written JSON, but there is no fsync guarantee, corruption recovery, persisted queue, or cross-process synchronization. Recovery updates are not written back; terminal save failure changes only memory ([fallback](../../server/main.go#L245)).
- **Confidence:** Mechanics confirmed; database-avoidance rationale inferred.

## Run FFmpeg behind a service-owned deadline

- **What:** The worker context derives from the service rather than the request; FFmpeg is invoked directly with restricted protocols/container formats, a two-minute deadline, first video and optional first audio, and VOD HLS output.
- **Evidence:** [worker context](../../server/main.go#L233), [command](../../server/main.go#L259).
- **Why, apparently:** Inferred: accepted work must outlive the 202 HTTP response while remaining cancelable at shutdown. Input restrictions and the timeout match the [documented limits](../../README.md#L20).
- **Tradeoff:** A browser disconnect after acceptance does not cancel processing; long/slow encodes fail after the fixed deadline. `ultrafast` favors encoding speed, and no resolution/bitrate ladder is configured.
- **Confidence:** Context/codec choices confirmed; performance rationale inferred from flags.

## Poll the collection and mount players only after completion

- **What:** Refresh the entire library every 1.5 seconds, rendering a player only for completed records with URLs.
- **Evidence:** [polling](../../client/src/App.tsx#L9), [conditional player](../../client/src/App.tsx#L25).
- **Why, apparently:** Inferred: polling gives simple eventual status updates without another transport or a job subscription API.
- **Tradeoff:** Status visibility has polling latency, collection cost grows with stored videos, and each completed item creates a player. The UI does not automatically retry rejected uploads.
- **Confidence:** Inferred from code; no explicit design note found.

## Player lifecycle belongs to a React effect

- **What:** Create the imperative Video.js element/player within an effect keyed by `src`, and call `dispose` during cleanup.
- **Evidence:** [VideoPlayer.tsx](../../client/src/VideoPlayer.tsx#L6), [StrictMode entry](../../client/src/main.tsx#L5).
- **Why, apparently:** Inferred: the effect provides a lifecycle boundary between React ownership and Video.js DOM/player ownership.
- **Tradeoff:** Source changes recreate the player rather than updating it in place. Cleanup is essential as components unmount or development StrictMode exercises effect setup/cleanup.
- **Confidence:** Lifecycle behavior confirmed; rationale inferred.

## Gotchas

- **Working directory controls storage.** `./videos` is resolved at launch, so starting the binary elsewhere gives a different library ([main.go:266](../../server/main.go#L266)).
- **Localhost is only the client's origin.** The UI hardcodes localhost, while the server binds all interfaces on `:8080`. A remotely opened UI would call the viewer's own machine unless reconfigured ([client origin](../../client/src/App.tsx#L4), [bind](../../server/main.go#L270)).
- **Collection order is not chronological.** The server sorts random IDs; upload time is not part of the model ([sort](../../server/main.go#L113), [model](../../server/main.go#L24)).
- **Failed conversion can leave partial output.** The worker deletes `input`, not its output directory. Restart recovery neither resumes processing nor cleans the input left by an abrupt stop ([worker](../../server/main.go#L244), [recovery](../../server/main.go#L67)).
- **Two-second HLS is a target.** The command sets `-hls_time 2` but no explicit keyframe interval; do not infer every segment has exactly that duration ([command](../../server/main.go#L259)).
- **A recovered network error can remain on screen.** Successful polling updates videos but does not clear the error state; beginning a new upload does clear it ([refresh](../../client/src/App.tsx#L12), [upload](../../client/src/App.tsx#L21)).
- **Old samples are not the current route.** The legacy `server/hls/` output does not match the new per-ID root/status gate ([README](../../README.md#L22), [media serving](../../server/main.go#L127)).
- **A test pass is not browser certification.** The encoder test conditionally skips without FFmpeg; with it, it checks files, not browser decoding. Earlier delivery verification records FFmpeg present and passing ([test](../../server/main_test.go#L54), [history](../../HISTORY.md#L56)).

## Conventions

- Keep the wire shape aligned between the Go `video` struct and the local TypeScript `Video` type; there is no shared schema generator ([server](../../server/main.go#L24), [client](../../client/src/App.tsx#L3)).
- Tests inject a replacement transcode function and use temporary directories, `httptest`, and explicit worker waits ([test](../../server/main_test.go#L16)). Follow that seam for deterministic lifecycle tests, then reserve real FFmpeg for integration coverage.
- Worker-facing detailed errors go to logs while clients receive stable generic failure messages ([error handling](../../server/main.go#L236)).

# Tech Stack

## Languages and runtimes

| Language or runtime | Version | Declared at |
| --- | --- | --- |
| Go | `1.23.5` module directive | [server/go.mod:3](../../server/go.mod#L3) |
| TypeScript | `~5.7.2`; browser target ES2020, Vite-config target ES2022 | [package.json:26](../../client/package.json#L26), [browser config](../../client/tsconfig.app.json#L4), [tool config](../../client/tsconfig.node.json#L4) |
| Node.js/npm | Runtime versions not pinned in the manifest | [package.json](../../client/package.json), [run instructions](../../README.md#L3) |
| Browser JavaScript | React application built as browser assets; no supported-browser matrix specified | [entry](../../client/src/main.tsx#L5), [browser caveat](../../README.md#L34) |
| FFmpeg | Executable version not pinned; requires libx264 and AAC encoding | [command](../../server/main.go#L259), [prerequisites](../../README.md#L3) |

## Frameworks and major libraries

Versions below are manifest ranges, not a claim about the exact installed lockfile resolution.

| Library | Version | Used for | Evidence |
| --- | --- | --- | --- |
| Go standard library | Tied to Go toolchain | HTTP, JSON, filesystem, random IDs, contexts, subprocesses, synchronization, signals. | [imports](../../server/main.go#L3), [Go version](../../server/go.mod#L3) |
| React and React DOM | Both `^19.0.0` | State/effects and DOM root rendering. | [manifest](../../client/package.json#L13), [App](../../client/src/App.tsx#L1), [entry](../../client/src/main.tsx#L1) |
| Video.js | `^8.21.0` | HLS source playback and player styling. | [manifest](../../client/package.json#L15), [player](../../client/src/VideoPlayer.tsx#L2) |
| Vite | `^6.1.0` | Development server and frontend asset build. | [manifest](../../client/package.json#L28), [scripts](../../client/package.json#L6) |
| Vite React plugin | `^4.3.4` | React integration for Vite. | [manifest](../../client/package.json#L21), [config](../../client/vite.config.ts#L2) |

## Data and infrastructure

| Service or facility | Role | Configured at |
| --- | --- | --- |
| Local filesystem | Per-ID input, atomic metadata replacement, generated HLS; no database. | [root](../../server/main.go#L266), [save](../../server/main.go#L76), [input](../../server/main.go#L211) |
| FFmpeg subprocess | Local conversion to H.264 video, optional AAC audio, VOD playlist and transport-stream segments. | [transcode](../../server/main.go#L258) |
| Go HTTP server | Port 8080; header timeout 5 seconds, read timeout 2 minutes, idle timeout 60 seconds. | [server configuration](../../server/main.go#L270) |
| In-process worker admission | Channel capacity one; two-minute encoder deadline; no broker. | [initialization](../../server/main.go#L49), [deadline](../../server/main.go#L233) |

The deployment description is explicitly a local, single-process demonstration, with no distributed job queue or production validation recorded ([README](../../README.md#L22), [history](../../HISTORY.md#L62)). No project-local container or CI configuration was found in the scoped file inventory; the [documented validation](../../README.md#L24) is command-based.

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| `go test -race` | Backend behavioral checks with concurrency instrumentation. | [README command](../../README.md#L28), [tests](../../server/main_test.go#L16) |
| npm scripts | `dev`, `build`, `lint`, and `preview`; build runs TypeScript then Vite. | [package.json:6](../../client/package.json#L6) |
| ESLint and `@eslint/js` | Both `^9.19.0`; recommended JavaScript baseline. | [manifest](../../client/package.json#L18), [config](../../client/eslint.config.js#L10) |
| `typescript-eslint` | `^8.22.0`; TypeScript lint configuration. | [manifest](../../client/package.json#L27), [config](../../client/eslint.config.js#L5) |
| React Hooks / React Refresh ESLint plugins | `^5.0.0` / `^0.4.18`; effect rules and component-export guidance. | [manifest](../../client/package.json#L23), [rules](../../client/eslint.config.js#L20) |
| React type packages | `@types/react ^19.0.8`, `@types/react-dom ^19.0.3`. | [manifest](../../client/package.json#L19) |

## Notes

- Backend dependencies are entirely standard library; FFmpeg is a runtime executable dependency outside Go's module graph ([go.mod](../../server/go.mod), [invocation](../../server/main.go#L259)).
- The frontend manifest has no test script. Existing automated behavioral coverage is backend-only; cross-platform browser playback is not certified ([scripts](../../client/package.json#L6), [README](../../README.md#L34)).
- [HISTORY.md](../../HISTORY.md#L54) records server race tests with real FFmpeg, frontend build, and lint passing for delivery. No new execution was performed to create this guide.

# Upload-to-HLS Project Guide

> Generated: 2026-09-21 from commit `2c4551f`. Scope: this `hls/` project only.

## What this is

This local demonstration accepts video uploads, converts them into HLS playlists and segments, and exposes processing status through a Go HTTP service ([server](../../server/main.go#L95)). A React library screen lets a user upload a file, poll its status, and play completed videos using Video.js ([UI](../../client/src/App.tsx#L9)). It consists of separate frontend and backend processes, an FFmpeg subprocess, and local files rather than a database or distributed queue ([implementation](../../server/main.go#L31)).

## Run it

From this project's root, start the backend:

```sh
cd server
go run .
```

In a second terminal, from this project's root:

```sh
cd client
npm ci
npm run dev
```

Open the URL printed by Vite. Go 1.23.5 is declared in [go.mod](../../server/go.mod#L3); Node.js/npm and FFmpeg with libx264/AAC support are also needed ([prerequisites](../../README.md#L3)). No environment variables are read: the UI uses `http://localhost:8080`, and the server binds `:8080` with storage at `./videos` relative to its working directory ([client](../../client/src/App.tsx#L4), [server](../../server/main.go#L266)). The bind address accepts connections on all interfaces, despite the root README describing localhost.

Documented checks, run from their respective folders:

```sh
# server/
go test -race ./...
# client/
npm run build
npm run lint
```

[HISTORY.md](../../HISTORY.md#L54) records earlier successful verification; this guide was produced by reading code, without rerunning builds or tests. The real FFmpeg test skips when FFmpeg is missing ([test](../../server/main_test.go#L54)).

## The five-file tour

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [server/main.go](../../server/main.go#L265) | Start at `main`, then trace upload admission, worker, and file serving. | [Startup](02-flow.md#startup) |
| 2 | [client/src/main.tsx](../../client/src/main.tsx#L5) | Mount the independent browser application. | [Components](01-architecture.md#components) |
| 3 | [client/src/App.tsx](../../client/src/App.tsx#L9) | Connect uploads and polling to server state. | [Upload and playback](02-flow.md#upload-and-playback) |
| 4 | [client/src/VideoPlayer.tsx](../../client/src/VideoPlayer.tsx#L4) | Translate a completed playlist URL into a disposable player. | [Player lifecycle](05-decisions.md#player-lifecycle-belongs-to-a-react-effect) |
| 5 | [server/main_test.go](../../server/main_test.go#L16) | See what upload, restart, and transcoding behavior is actually verified. | [Verification](02-flow.md#build-and-verification) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, contracts, and persistence.
2. [Flow](02-flow.md) — startup through playback and shutdown.
3. [Structure](03-structure.md) — files and callers.
4. [Tech stack](04-tech-stack.md) — versions and tooling.
5. [Decisions](05-decisions.md) — evidence, tradeoffs, and gotchas.

## Open questions

- What production origin, bind address, authentication, retention, and storage quota are intended? The current configuration is hardcoded and the README explicitly scopes this to a local demonstration ([configuration](../../server/main.go#L266), [limitations](../../README.md#L22)).
- Which Node.js and FFmpeg versions and browser playback combinations are supported? The [manifest](../../client/package.json#L1) does not pin a Node runtime, the [FFmpeg call](../../server/main.go#L259) selects the executable from PATH, and [browser certification](../../README.md#L34) remains open.
- Should restart recovery persist the transition to `failed` or clean interrupted inputs? The [loader](../../server/main.go#L67) changes only the in-memory record.
- How should busy admission, cancellation, persistence failure, and browser interactions be verified? [V1.md](../../V1.md#L9) calls for lifecycle failure tests, but the current [three tests](../../server/main_test.go#L16) cover successful upload/reload, real transcoding, and an invalid multipart request.

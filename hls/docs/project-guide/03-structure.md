# Structure

## What lives where

The project has two independently launched applications. [server/](../../server/) contains the Go service and its behavioral tests; [client/src/](../../client/src/) contains the browser bootstrap, library screen, and player adapter. The [root README](../../README.md), [v1 contract](../../V1.md), and [history](../../HISTORY.md) document usage, scope, and prior delivery evidence.

```text
hls/
├── README.md, V1.md, HISTORY.md
├── server/
│   ├── go.mod
│   ├── main.go
│   ├── main_test.go
│   ├── hls/                  old generated sample output
│   └── videos/               runtime storage, created on startup
├── client/
│   ├── package.json
│   ├── index.html
│   ├── vite.config.ts
│   ├── eslint.config.js
│   ├── tsconfig*.json
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       └── VideoPlayer.tsx
└── docs/project-guide/       this guide
```

The displayed runtime storage location assumes `go run .` from `server/`; [main](../../server/main.go#L266) resolves it from the process working directory.

## Server

| File | Responsibility | Key exports or symbols | Called by |
| --- | --- | --- | --- |
| [server/main.go](../../server/main.go) | JSON video model, metadata restore/save, route switch, admission, transcoder, and process lifecycle. | `video`, `service`, `newService`, `validID`, `save`, `Close`, `ServeHTTP`, `upload`, `transcode`, `main` (package-local types/functions except methods). | Go entry point; `http.Server`; tests; worker goroutine. |
| [server/main_test.go](../../server/main_test.go) | Upload/media/restart test, optional real encoder test, invalid multipart test. | `TestUploadProcessingAndRestart`, `TestFFmpegPipeline`, `TestRejectInvalidUpload`. | `go test`. |
| [server/go.mod](../../server/go.mod) | Declares standalone `server` module and Go version; no external module requirements. | Module metadata. | Go tooling. |

All server responsibilities deliberately remain in one source file; use the [flow](02-flow.md) anchors to navigate it.

## Client source

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [client/src/main.tsx](../../client/src/main.tsx) | Mount `App` into the HTML root under StrictMode. | Side-effect entry module. | HTML module script through Vite. |
| [client/src/App.tsx](../../client/src/App.tsx) | Hold videos/error/upload state, POST files, poll the library, render completed players. | Default `App`; local `Video` type and `api` constant. | Browser entry. |
| [client/src/VideoPlayer.tsx](../../client/src/VideoPlayer.tsx) | Own imperative Video.js initialization and disposal; load Video.js CSS. | Default `VideoPlayer({ src })`. | `App` for each completed video with a URL. |

## Client setup

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [client/index.html](../../client/index.html) | HTML root and module entry; retains template page title/favicon. | `#root` mount point. | Vite/browser. |
| [client/package.json](../../client/package.json) | Runtime/development dependencies and dev/build/lint/preview scripts. | npm scripts. | npm tooling. |
| [client/vite.config.ts](../../client/vite.config.ts) | Enable the React Vite plugin; no API proxy or origin configuration. | Default Vite config. | Vite. |
| [client/eslint.config.js](../../client/eslint.config.js) | Lint TS/TSX, ignore dist, and configure Hooks/Refresh rules. | Default flat config. | ESLint. |
| [client/tsconfig.json](../../client/tsconfig.json) | Reference browser and Vite-tool projects for `tsc -b`. | Project references. | TypeScript build. |
| [client/tsconfig.app.json](../../client/tsconfig.app.json) | Strict browser/React compilation, ES2020 target, bundler resolution, no emitted JS. | Compiler options for `src`. | TypeScript build. |
| [client/tsconfig.node.json](../../client/tsconfig.node.json) | Strict Vite-config compilation with ES2022 target and ES2023 libraries. | Compiler options for Vite config. | TypeScript build. |

## Project documents

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [README.md](../../README.md) | Current run commands, endpoints, validation commands, and local-demo limits. | None. | Maintainers/users. |
| [V1.md](../../V1.md) | Upload-to-HLS acceptance contract. | None. | Delivery/review work. |
| [HISTORY.md](../../HISTORY.md) | Git-backed chronology and recorded delivery verification. | None. | Maintainers. |
| [client/README.md](../../client/README.md) | Retained Vite template guidance, not the current product's operation guide. | None. | Frontend maintainers. |

## Excluded

Generated sample segments and playlist in `server/hls/` are excluded: the [root README](../../README.md#L22) says this earlier fixed sample is no longer the upload workflow. Runtime `videos/`, installed dependencies, build output, dependency lockfiles, template SVG assets, and the one-line Vite type-reference boilerplate are omitted from the source tables. Their omission does not imply they are production components; [the media route](../../server/main.go#L145) serves only the configured runtime root.

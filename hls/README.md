# Upload-to-HLS v1

Requires Go, Node.js and FFmpeg with libx264/AAC support.

```sh
cd server
go run .
# In another terminal, from hls/client:
npm ci
npm run dev
```

Open the client URL printed by Vite. Upload a video, watch its processing state, and play its generated HLS playlist. The server listens on localhost port 8080; the development UI uses that endpoint.

- POST `/api/videos`: multipart `file`, 1 byte to 100 MiB; returns 202 with an ID and processing state.
- GET `/api/videos`: current library.
- GET `/api/videos/<id>`: one processing/completed/failed status.
- GET `/hls/<id>/index.m3u8` and segment URLs: completed video output.

One processor runs at a time; busy uploads return 503 so the user can retry. Each upload has a random private storage directory. Input filenames are display labels, not storage paths. FFmpeg has a two-minute deadline and accepts video containers rather than playlist/concat inputs. Optional audio is retained when available. Original inputs are removed after processing. Metadata survives restart; interrupted jobs become failed. Disk files are atomically replaced but this service does not promise crash-durable metadata fsync.

SIGINT/SIGTERM stops admission, cancels processing and drains HTTP work with a shutdown deadline. V1 is a local single-process demonstration: storage has no automatic quota/retention policy, authentication or distributed job queue. The earlier fixed `hls/output.m3u8` sample is no longer the upload workflow.

Validation:

```sh
# server
go test -race ./...
# client
npm run build
npm run lint
```

The backend suite generates a short input with FFmpeg and verifies a finished playlist and transport-stream output, plus upload/status/restart behavior. Browser playback across every platform has not been certified.

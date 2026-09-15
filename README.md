# Go Core Implementations

A collection of educational implementations and experiments in Go, covering data structures, concurrency, networking, storage, and distributed systems. Projects are independent and vary in completeness. Most have their own `go.mod`; the HLS project also includes a React/TypeScript client. Additional projects are linked as Git submodules.

## Getting started

Clone with submodules to fetch all projects:

```bash
git clone --recurse-submodules https://github.com/rushikeshg25/go-core-impl.git
cd go-core-impl
```

For an existing checkout, initialize the submodules with:

```bash
git submodule update --init --recursive
```

Use the Go version declared in each project's `go.mod` or a newer compatible version. The modules stored directly in this repository currently declare versions from Go 1.21 through Go 1.25.0. Run Go commands from the relevant module directory; there is no root Go module.

For example:

```bash
cd consistent-hashing
go run .
```

Some projects require MySQL, Redis, Kafka, or a JavaScript runtime. Check their source configuration and linked documentation before running them. Several network demos default to port `8080`, so run them separately or adjust their ports.

## Projects in this repository

### [Adding NULL vs NOT NULL Column Benchmarking](adding-null-vs-not-null-col-benchmarking/)

Compares MySQL `ALTER TABLE` timings for adding a nullable integer column versus a `NOT NULL DEFAULT 0` column, using 100,000 rows over five iterations. It also reports Go process memory statistics. Requires MySQL and a `test` database; the benchmark creates and drops the `alter_benchmark` table.

### [Bloom Filter](bloomfilter/)

A Bloom filter library using seeded MurmurHash3 hashes, with `New`, `Add`, `Test`, and `Clear` operations. Supports probabilistic membership checks with configurable storage size and hash count. This is a library package, with no executable demo.

### [Concurrency Control](concurrency-control/)

A Fiber HTTP API backed by MySQL that demonstrates optimistic concurrency control using either version numbers or SHA-256 checksums. Clients submit a conflict token when updating a record, and stale updates receive HTTP `409`. Includes tests and a MySQL Docker Compose configuration. Run with `go run . -v` for version mode or `go run . -c` for checksum mode after starting MySQL.

### [Consistent Hashing](consistent-hashing/)

A hash ring using SHA-256 and sorted node positions. Supports adding and removing storage nodes, assigning keys to nodes, and printing the ring's positions. Includes a small command-line demo.

### [Durable Logs](durable-logs/)

An experimental file logger that serializes timestamped entries with Protocol Buffers, buffers writes, and rotates log segments after a configured entry count. Includes a logging demo; reading buffered logs is currently a placeholder.

### [Event Loop](event-loop/)

An experimental event loop with separate channels for tasks and callbacks, plus a bounded pool of goroutines for asynchronous tasks. The source defines the loop API but currently has no `main` entry point.

### [HLS (HTTP Live Streaming)](hls/)

A Go HTTP server for file uploads and serving existing HLS playlists and segments, paired with a React/TypeScript client using Video.js. The server does not currently transcode uploads or generate playlists. The server module is in `hls/server`; the Vite client and its package scripts are in `hls/client`.

### [Kafka Consumer Groups](kafka-multiple-consumers-partitions/)

Command-line examples of partitioned producers, consumer groups with manual offset commits, and a transactional producer using Confluent's Kafka client. Also includes printed walkthroughs of dead-letter queues, sagas, monitoring, and schema evolution. Messaging examples expect Kafka at `localhost:9092`; run without arguments to list the available modes.

### [Merkle Trees](merkle-trees/)

A directory Merkle tree CLI using BLAKE3 hashes. Supports `hash <dir>` to print a root hash, `print <dir>` to display the tree, and `diff <dir1> <dir2>` to report added, deleted, or modified entries. For example, run `go run . hash .` from the module directory.

### [Mini Git](mini-git/)

A small version-control CLI built with Cobra, with commands for `init`, `add`, `commit`, `status`, `log`, `branch`, and `checkout`. Includes tests and a Makefile for building the executable. Run `go run ./cmd --help` from the module directory to see the commands.

### [Multithreaded TCP Server](multithreaded-tcp/)

A TCP broadcast server using a goroutine per connection, a buffered message channel, and a mutex-protected client registry. Reads newline-delimited messages and broadcasts them to connected clients. Listens on port `8080` and includes signal-driven shutdown handling.

### [Queue](queue/)

A slice-backed integer FIFO queue with enqueue and dequeue operations that accept a caller-supplied mutex. Includes a demo that launches concurrent producers and consumers.

### [RAFT](RAFT/)

An in-progress Raft implementation demonstrating leader election, randomized timeouts, and heartbeats over Go's `net/rpc`. Log replication and persistent storage remain planned. See the [RAFT README](RAFT/README.md) for the three-node demo and roadmap.

### [Real-time Leaderboard](realtime-leaderboard/)

A Redis sorted-set leaderboard exposed through an HTTP API for updating scores, retrieving a player's rank, and listing the top ten players. Includes tests and environment-based configuration. See the [leaderboard README](realtime-leaderboard/README.md) for setup and API examples.

### [Thread Pool](thread-pool/)

Two worker-pool demos using goroutines and channels: a task-function pool in `main.go` and a jobs/results-channel example in `approach1.go`. Each file has its own `main` function, so run them separately with `go run main.go` or `go run approach1.go`.

### [WebSockets](websockets/)

A manual WebSocket HTTP upgrade handshake, including computation of `Sec-WebSocket-Accept`, plus a health endpoint. The connection closes after the handshake; WebSocket frame exchange and persistent messaging are not implemented yet.

## Projects linked as Git submodules

These projects live in separate repositories, pinned to commits by this repository. Their directories may be empty until submodules are initialized. Follow the upstream links for implementation details and setup instructions.

| Project | Directory | Upstream repository |
| --- | --- | --- |
| Adaptive Bitrate Streaming | [adaptive-bitrate-streaming](adaptive-bitrate-streaming/) | [Source](https://github.com/rushikeshg25/adaptive-bitrate-streaming) |
| B+ Tree | [bp-tree](bp-tree/) | [Source](https://github.com/rushikeshg25/bp-tree) |
| Load Balancer | [loadbalancer](loadbalancer/) | [Source](https://github.com/rushikeshg25/loadbalancer) |
| P2P File Sharing | [p2p-file-sharing](p2p-file-sharing/) | [Source](https://github.com/rushikeshg25/p2p-file-sharing) |
| Task Scheduler | [task-scheduler](task-scheduler/) | [Source](https://github.com/rushikeshg25/task-scheduler) |
| Token Bucket | [token-bucket](token-bucket/) | [Source](https://github.com/rushikeshg25/token-bucket) |
| Tricolor Garbage Collection | [tricolor-gc](tricolor-gc/) | [Source](https://github.com/rushikeshg25/tricolor-gc) |
| WAL (Write-Ahead Log) | [wal-go](wal-go/) | [Source](https://github.com/rushikeshg25/wal-go) |

## Keeping this index current

When adding, renaming, or removing a project, update its entry and link here in the same change. Describe the behavior currently implemented, identify unfinished features, and document any special entry point or external service requirement. For submodules, keep the directory and upstream URL aligned with [.gitmodules](.gitmodules).

# Go Core Implementations

Implementations of core systems concepts in Go, spanning data structures, concurrency, storage, networking, and distributed systems.

This repository brings together independent libraries, command-line tools, and service demos for studying how these systems work. Each project focuses on a specific concept, with its own dependencies and development status. The HLS example also includes a React and TypeScript client.

[Projects](#projects) · [Git Submodules](#git-submodules) · [Getting Started](#getting-started) · [Development](#development)

## Projects

### Data Structures and Storage

| Project | Implementation |
| --- | --- |
| [Bloom Filter](bloomfilter/) | Probabilistic membership testing with seeded MurmurHash3 hashes and configurable storage size and hash count. Exposes `New`, `Add`, `Test`, and `Clear` as a library. |
| [Durable Logs](durable-logs/) | Buffered Protocol Buffer logging with exact framing, count-based rotation, pending-log retrieval, restart recovery, and explicit sync/close errors. |
| [Merkle Trees](merkle-trees/) | Directory hashing and comparison using BLAKE3. Provides commands to calculate root hashes, display trees, and report added, deleted, or modified entries. |
| [Mini Git](mini-git/) | A Cobra-based version-control CLI with repository initialization, staging, commits, status, history, branches, and checkout. |
| [Queue](queue/) | A slice-backed integer FIFO queue with a caller-supplied mutex and a concurrent producer/consumer demo. |

### Concurrency and Database Experiments

| Project | Implementation |
| --- | --- |
| [Concurrency Control](concurrency-control/) | A MySQL-backed Fiber API demonstrating optimistic concurrency control through version numbers or SHA-256 checksums. Rejects stale conflict tokens with HTTP `409`. |
| [Event Loop](event-loop/) | Runnable task/callback loop with bounded asynchronous execution, cooperative per-event cancellation, callback completion, and draining shutdown. |
| [Schema Change Benchmark](adding-null-vs-not-null-col-benchmarking/) | Compares MySQL column additions using `NULL` and `NOT NULL DEFAULT 0` across five iterations of 100,000 rows. Reports execution time and Go process memory statistics. |
| [Thread Pool](thread-pool/) | Two worker-pool examples: a task-function pool and a jobs/results-channel implementation. Each runs as a separate program. |

### Distributed Systems

| Project | Implementation |
| --- | --- |
| [Consistent Hashing](consistent-hashing/) | A SHA-256 hash ring with node addition, removal, key assignment, and a command-line visualization of node positions. |
| [Kafka Consumer Groups](kafka-multiple-consumers-partitions/) | Partitioned producers, consumer groups with manual offset commits, and transactional production using Confluent's Kafka client. Includes conceptual walkthroughs of messaging patterns and operations. |
| [RAFT](RAFT/) | Persistent static-membership consensus with election, replicated logs, conflict repair, majority commit, ordered application checkpoints, and bounded TCP/RPC. See the [demo and roadmap](RAFT/README.md). |

### Storage Engine Series

These build up a storage engine one concept at a time, each in its own repository. They are independent Go modules with no dependency on one another: shared ideas are re-implemented rather than factored out, so each can be read on its own. All are tracked here as [submodules](#git-submodules).

| Project | Implementation |
| --- | --- |
| [Append-Only Storage](append-only-storage/) | Versioned record framing, positional reads and scans, incomplete-tail recovery, explicit sync, and read-only inspection. See its v1 contract and project guide. |
| [Checksum and Corruption Detection](checksum-corruption-detection/) | A CRC32C record log with a read-only verifier and a `crcverify` command. Reports the offset and both checksums of a damaged record, and repairs only a final record that EOF cut short. |
| [In-Memory Hash Index](in-memory-hash-index/) | An on-disk key-value store whose live keys resolve through a custom open-addressed hash table with linear probing, deleted-slot tombstones, and probe-distribution statistics. |
| [Concurrent Readers, Single Writer](concurrent-readers-single-writer/) | The same engine with concurrent reads and serialized mutations, a documented lock ordering, and deterministic tests over operation publication schedules. |
| [Tombstone Deletes](tombstone-deletes/) | Logical deletion by appending tombstone records, with delete misses that write nothing, resurrection by a later put, and statistics for the space obsolete records hold. |
| [Storage Compaction](storage-compaction/) | Reclaims that space by rewriting only live records into a candidate file and publishing it with an atomic rename, with fault injection at every protocol phase and subprocess crash tests. |

### Networking and Applications

| Project | Implementation |
| --- | --- |
| [HTTP Live Streaming](hls/) | Isolated asynchronous uploads and FFmpeg-to-HLS conversion with persisted processing status and a connected React/TypeScript Video.js player. |
| [Multithreaded TCP Server](multithreaded-tcp/) | A newline-delimited TCP broadcast server with a goroutine per connection, buffered message delivery, and signal-driven shutdown handling. |
| [Real-time Leaderboard](realtime-leaderboard/) | A Redis sorted-set HTTP API for score updates, player ranks, and the top ten players. Includes tests and environment-based configuration. See the [setup and API reference](realtime-leaderboard/README.md). |
| [WebSockets](websockets/) | Manual RFC 6455 echo with validated upgrade, masked frames, fragmentation, text/binary messages, ping/pong, and close handling. |

## Git Submodules

The following projects are maintained in separate repositories and pinned to specific commits. Initialize the submodules to populate their local directories. The storage engine projects are described in the [catalog above](#storage-engine-series); for the rest, refer to each upstream repository for its implementation details and setup requirements.

| Project | Local Directory | Upstream |
| --- | --- | --- |
| Adaptive Bitrate Streaming | [adaptive-bitrate-streaming](adaptive-bitrate-streaming/) | [Repository](https://github.com/rushikeshg25/adaptive-bitrate-streaming) |
| Append-Only Storage | [append-only-storage](append-only-storage/) | [Repository](https://github.com/rushikeshg25/append-only-storage) |
| B+ Tree | [bp-tree](bp-tree/) | [Repository](https://github.com/rushikeshg25/bp-tree) |
| Checksum and Corruption Detection | [checksum-corruption-detection](checksum-corruption-detection/) | [Repository](https://github.com/rushikeshg25/checksum-corruption-detection) |
| Concurrent Readers, Single Writer | [concurrent-readers-single-writer](concurrent-readers-single-writer/) | [Repository](https://github.com/rushikeshg25/concurrent-readers-single-writer) |
| In-Memory Hash Index | [in-memory-hash-index](in-memory-hash-index/) | [Repository](https://github.com/rushikeshg25/in-memory-hash-index) |
| Load Balancer | [loadbalancer](loadbalancer/) | [Repository](https://github.com/rushikeshg25/loadbalancer) |
| Monsoon | [monsoon](monsoon/) | [Repository](https://github.com/rushikeshg25/monsoon) |
| P2P File Sharing | [p2p-file-sharing](p2p-file-sharing/) | [Repository](https://github.com/rushikeshg25/p2p-file-sharing) |
| Storage Compaction | [storage-compaction](storage-compaction/) | [Repository](https://github.com/rushikeshg25/storage-compaction) |
| Task Scheduler | [task-scheduler](task-scheduler/) | [Repository](https://github.com/rushikeshg25/task-scheduler) |
| Token Bucket | [token-bucket](token-bucket/) | [Repository](https://github.com/rushikeshg25/token-bucket) |
| Tombstone Deletes | [tombstone-deletes](tombstone-deletes/) | [Repository](https://github.com/rushikeshg25/tombstone-deletes) |
| Tricolor Garbage Collection | [tricolor-gc](tricolor-gc/) | [Repository](https://github.com/rushikeshg25/tricolor-gc) |
| Write-Ahead Log | [wal-go](wal-go/) | [Repository](https://github.com/rushikeshg25/wal-go) |

## Getting Started

### Clone the Repository

```bash
git clone --recurse-submodules https://github.com/rushikeshg25/go-core-impl.git
cd go-core-impl
```

To initialize submodules in an existing checkout:

```bash
git submodule update --init --recursive
```

### Prerequisites

Use the Go version specified in the selected project's `go.mod`, or a newer compatible version. The modules maintained directly in this repository currently declare versions from Go 1.21 to Go 1.25.0.

Additional requirements depend on the project:

| Project | Requirements |
| --- | --- |
| Concurrency Control | MySQL. A Docker Compose configuration is included in the project directory. |
| Schema Change Benchmark | MySQL with a `test` database. The benchmark creates and drops the `alter_benchmark` table. |
| Kafka Consumer Groups | Kafka at `localhost:9092` for the messaging examples. |
| Real-time Leaderboard | Redis. Connection settings are configurable through environment variables. |
| HTTP Live Streaming | FFmpeg for server transcoding, plus a JavaScript runtime and package manager for the Vite client. |

Several network demos use port `8080` by default. Run them separately or adjust their port configuration.

### Run a Project

Run commands from the relevant module directory. For example:

```bash
cd consistent-hashing
go run .
```

Projects with specific entry points or arguments use the following commands, relative to their module directories:

| Project | Command |
| --- | --- |
| Checksum and Corruption Detection | `go run ./cmd/crcverify [-json] <file>` to check a log for damage |
| Concurrency Control | `go run . -v` for version mode or `go run . -c` for checksum mode |
| Merkle Trees | `go run . hash <dir>`, `go run . print <dir>`, or `go run . diff <dir1> <dir2>` |
| Mini Git | `go run ./cmd --help` |
| Thread Pool | `go run main.go` or `go run approach1.go` |

For Kafka, running without arguments lists the available modes. For RAFT, follow the three-node instructions in its README. The HLS server module is in `hls/server`, and the client package scripts are in `hls/client/package.json`.

## Development

Each project is maintained independently. There is no root Go module or repository-wide build command. Build and test within the selected module, using its documented entry point and required services. Bloom Filter is a library; Event Loop ships a runnable task/callback demo. The storage engine projects are libraries as well, each with its own test suite; Checksum and Corruption Detection also ships the `crcverify` command.

When adding or changing a project, update this catalog in the same change. Keep descriptions aligned with implemented behavior, record relevant setup requirements, and keep submodule paths and upstream URLs consistent with [.gitmodules](.gitmodules).

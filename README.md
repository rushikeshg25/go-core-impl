# Go Core Implementations

This repository serves as a comprehensive collection of low-level system implementations, distributed algorithms, and core computer science concepts authored in Go. Each subdirectory contains a self-contained module demonstrating a specific technical challenge or architectural pattern.

## Table of Contents

- [Core Data Structures and Algorithms](#core-data-structures-and-algorithms)
- [Distributed Systems](#distributed-systems)
- [Networking and Communication](#networking-and-communication)
- [Concurrency and Performance](#concurrency-and-performance)
- [Storage and Logging](#storage-and-logging)
- [Streaming and Multimedia](#streaming-and-multimedia)
- [External Modules (Submodules)](#external-modules-submodules)
- [Usage Information](#usage-information)

---

## Core Data Structures and Algorithms

### Bloom Filter
Location: `bloomfilter/`
A probabilistic data structure used to test whether an element is a member of a set. It provides a highly space-efficient way to handle membership queries, with a possibility of false positives but no false negatives.
- Implementation: Uses MurmurHash3 for efficient hashing.
- Features: Configurable size and number of hash functions (k).

### Event Loop
Location: `event-loop/`
An implementation of an asynchronous execution model, similar to those found in Node.js or browser environments.
- Components: Event queue, callback queue, and a central loop.
- Capabilities: Supports both synchronous and asynchronous task execution with callback management.

### Real-time Leaderboard
Location: `realtime-leaderboard/`
A high-performance leaderboard system designed for real-time updates and queries.
- Backend: Built using Redis for atomic operations and low-latency data access.
- Functionality: Supports score updates, rank retrieval, and top-k queries.

---

## Distributed Systems

### Consistent Hashing
Location: `consistent-hashing/`
A distributed hashing scheme used to minimize reorganization when nodes are added or removed from a cluster.
- Implementation: Uses a hash ring to map both items and storage nodes.
- Advantages: Ensures data distribution stability, making it ideal for load balancers and distributed caches.

### Kafka Consumer Groups and Advanced Patterns
Location: `kafka-multiple-consumers-partitions/`
A detailed exploration of Apache Kafka's advanced features and architectural patterns.
- Consumer Groups: Demonstrates how partitions are balanced across multiple consumers for horizontal scaling.
- Exact-Once Semantics: Implementation of transactional producers to guarantee message delivery exactly once.
- Advanced Patterns: Includes documentation and conceptual code for Dead Letter Queues (DLQ) and the Saga pattern for distributed transactions.

---

## Networking and Communication

### Multithreaded TCP Server
Location: `multithreaded-tcp/`
A concurrent TCP server implementation capable of handling multiple simultaneous client connections.
- Model: Each incoming connection is handled by a separate goroutine, demonstrating Go's lightweight concurrency primitives.
- Features: Basic request-response handling over raw TCP sockets.

### WebSockets
Location: `websockets/`
Implementations of full-duplex communication channels over a single TCP connection.
- Use Case: Real-time interactive applications requiring low-latency bi-directional messaging.

---

## Concurrency and Performance

### Thread Pool
Location: `thread-pool/`
A resource management pattern that maintains a pool of available workers to execute tasks concurrently.
- Implementation: Uses Go channels for task distribution and WaitGroups for synchronization.
- Objective: Limits the overhead of frequent thread/goroutine creation and destruction by reusing existing workers.

---

## Storage and Logging

### Durable Logs
Location: `durable-logs/`
A system for persistent message logging designed for reliability and crash recovery.
- Serialization: Uses Protocol Buffers (Protobuf) for structured, efficient data storage.
- Storage: Implements disk-backed logging to ensure data survives process restarts.

---

## Streaming and Multimedia

### HTTP Live Streaming (HLS)
Location: `hls/`
A comprehensive implementation of the HLS protocol for video delivery.
- Server: Handles segmenting and playlist (m3u8) generation.
- Client: Capable of consuming HLS streams and managing buffer state.

---

## External Modules (Submodules)

This repository integrates several standalone core projects as Git submodules:

- **Write-Ahead Log (WAL)** (`wal-go`): A foundational component for database durability and crash recovery.
- **Load Balancer** (`loadbalancer`): Implementations of various traffic distribution algorithms.
- **Tricolor Garbage Collection** (`tricolor-gc`): A demonstration of the tricolor marking algorithm used in modern GCs.
- **Token Bucket** (`token-bucket`): A rate-limiting algorithm for controlling network traffic and API usage.
- **P2P File Sharing** (`p2p-file-sharing`): Peer-to-peer communication and data transfer protocols.
- **B+ Tree** (`bp-tree`): A self-balancing tree data structure widely used in database indexing.
- **Adaptive Bitrate Streaming** (`adaptive-bitrate-streaming`): Logic for dynamically adjusting video quality based on network conditions.

---

## Usage Information

### Prerequisites
- Go 1.21 or higher.
- Redis (required for `realtime-leaderboard`).
- Kafka and Zookeeper (required for `kafka-multiple-consumers-partitions`).

### Building and Running
Most sub-projects can be executed directly using `go run .` within their respective directories. A global `makefile` is also provided for convenience:

```bash
# Example: Run consistent hashing demonstration
make consistent-hashing

# Example: Run thread pool demonstration
make thread-pool
```

### Note on Submodules
If you have just cloned the repository, ensure submodules are initialized:
```bash
git submodule update --init --recursive
```

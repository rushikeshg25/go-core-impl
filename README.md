# Go Core Implementations

This repository serves as a comprehensive collection of low-level system implementations, distributed algorithms, and core computer science concepts authored in Go. Each subdirectory contains a self-contained module demonstrating a specific technical challenge or architectural pattern.

### Adaptive Bitrate Streaming
Implments logic for dynamically adjusting video quality based on network conditions. It monitors bandwidth availability and switches between different stream profiles to ensure smooth playback. This project demonstrates real-time decision-making in streaming media.

### Adding NULL vs NOT NULL Column Benchmarking
A specialized benchmarking tool designed to measure the performance impact of database schema changes. It specifically compares the speed and resource usage of adding NULLable versus NOT NULL columns. This helps in understanding low-level database engine behavior during migrations.

### Bloom Filter
A probabilistic, space-efficient data structure used for rapid membership testing. It uses MurmurHash3 to minimize false positives while maintaining zero false negatives. This implementation is ideal for caching layers and large-scale set operations where memory is constrained.

### B+ Tree
A self-balancing tree data structure widely used in modern database indexing and file systems. It ensures efficient data storage, retrieval, and range queries by maintaining sorted data and providing logarithmic search time. This project showcases complex pointer management and node splitting/merging logic.

### Consistent Hashing
A distributed hashing scheme that provides a stable way to map data across a dynamic cluster of nodes. It minimizes data reorganization when servers are added or removed, making it perfect for load balancers and distributed caches. The implementation uses a hash ring for efficient resource allocation.

### Durable Logs
A persistent message logging system engineered for high reliability and crash recovery. It utilizes Protocol Buffers for structured data serialization and implements disk-backed storage to ensure durability. This project is a core component for building databases and message brokers.

### Event Loop
An implementation of an asynchronous execution model similar to those found in Node.js. It features a central loop, task queues, and callback management to enable non-blocking I/O operations. This module demonstrates how to manage concurrency without traditional threading overhead.

### HLS (HTTP Live Streaming)
A complete implementation of the Apple HLS protocol for modern video delivery. It handles segmenting raw video files and generating m3u8 playlists for adaptive streaming. This server-side implementation supports real-time playback across various devices and network conditions.

### Kafka Consumer Groups
An exploration of advanced Apache Kafka architectural patterns and features. It demonstrates partition balancing across consumer groups, exact-once delivery semantics, and transactional producer logic. This is a practical guide for building scalable, fault-tolerant distributed streaming applications.

### Load Balancer
Implementing various traffic distribution algorithms, this project enhances application availability and performance. It explores strategies like Round Robin and Least Connections to manage requests across a backend cluster. This serves as a foundational tool for scaling web services.

### Multithreaded TCP Server
A concurrent TCP server capable of handling thousands of simultaneous client connections. It leverages Go's lightweight goroutines to process each connection independently without blocking the main thread. This project demonstrates the power of Go's concurrency primitives for network programming.

### P2P File Sharing
Implementation of peer-to-peer communication protocols for decentralized data transfer. It explores neighbor discovery, chunk-based file distribution, and data integrity verification in a distributed network. This project provides insight into the architecture of systems like BitTorrent.

### Queue
A collection of fundamental queue implementations for managing data flow and task execution. It includes thread-safe operations and demonstrates various queuing strategies for different use cases. These components are essential building blocks for any concurrent system.

### Real-time Leaderboard
A high-performance leaderboard system engineered for low-latency score updates and ranking. It utilizes Redis for atomic operations and fast data retrieval, supporting top-k queries and rank tracking. This implementation is designed for gaming and social applications with high throughput.

### Task Scheduler
A robust system for scheduling and managing background tasks with integrated monitoring. It features a visual terminal UI for real-time status tracking and job management. This project demonstrates complex task orchestration and interactive CLI design.

### Thread Pool
A resource management pattern that maintains a set of workers to execute tasks concurrently. It optimizes performance by reducing the overhead of constant goroutine creation and destruction. This implementation uses Go channels for efficient task distribution and synchronization.

### Token Bucket
A classic rate-limiting algorithm used for traffic shaping and controlling API usage. It allows for bursts of traffic while maintaining a steady average rate, protecting backend services from overload. This module is a key component for implementing fair-use policies in web applications.

### Tricolor Garbage Collection
A demonstration of the tricolor marking algorithm used in modern garbage collectors. It visualizes the process of identifying reachable objects and reclaiming memory in a concurrent environment. This project provides deep insight into the internal workings of Go's runtime.

### WAL (Write-Ahead Log)
A foundational component for database durability, ensuring all data changes are logged before being applied. It provides the necessary primitives for crash recovery and transaction atomicity. This implementation is critical for anyone building custom storage engines or databases.

### WebSockets
Implementation of full-duplex communication channels over a single TCP connection. It enables real-time, bi-directional messaging between clients and servers for interactive applications. This project demonstrates handling persistent connections and low-latency data exchange.


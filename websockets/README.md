# WebSocket echo server v1

Run `go run .`, then connect to `ws://localhost:8080/ws`. `/health` returns a simple health response. Text and binary messages are echoed; fragmented messages are assembled, with ping/pong processing between fragments.

The HTTP upgrade validates GET, Connection/Upgrade tokens, version 13 and a 16-byte decoded key. Browser Origin headers must match the request host. Clients without Origin are accepted. The server preserves bytes buffered during upgrade.

Clients must mask frames. Reserved bits/opcodes, noncanonical lengths, fragmented control frames and invalid close payloads are rejected. Text messages require valid UTF-8. Frames and assembled messages are bounded to 1 MiB; control payloads to 125 bytes. Close requests are echoed; invalid input receives a close frame before disconnection. Extensions, compression and subprotocol negotiation are outside v1.

Idle reads time out after 60 seconds, writes after 10 seconds. No authentication or TLS is included; deploy behind a suitable endpoint when needed. Process shutdown closes connections; this demo does not track hijacked connections for a graceful server drain.

Run `go test -race ./...` for raw HTTP/TCP upgrade, fragmented echo, ping/pong, close and malformed-frame checks. `go test -run '^$' -fuzz FuzzFrame -fuzztime 5s` exercises the frame parser.

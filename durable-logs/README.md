# Durable logs v1

A synchronized buffered protobuf logger with numbered segments and final-tail recovery.

```go
d, err := durablelogs.Open("./logs", 100) // records per segment
if err != nil { panic(err) }
if err := d.Log("hello"); err != nil { panic(err) }
pending := d.GetBufferedLogs()
if err := d.Flush(); err != nil { panic(err) }
records, err := d.ReadAll()
if err := d.Close(); err != nil { panic(err) }
```

`GetBufferedLogs` returns an owned snapshot of messages accepted since the last successful flush, including any bytes automatically drained by bufio. Flush drains bufio, fsyncs the current file and clears that pending list. ReadAll flushes before reading every segment in numeric order. Rotation occurs before the next record when the current segment reaches its configured count.

The existing format is preserved: little-endian uint32 length followed by a protobuf Log containing message and timestamp. Encoded records are bounded to 64 MiB. Open validates every segment and repairs only an incomplete record at the end of the latest segment. Earlier incomplete records, invalid protobuf and oversized lengths fail open. There are no checksums, so arbitrary bit corruption is not reliably detected.

All public operations synchronize on one mutex. The first write/flush failure poisons the instance; subsequent mutations return ErrPoisoned until reopen. Close attempts flush and file close, joins errors, and is idempotent. Other operations that return errors use ErrClosed afterward. One process must exclusively own the directory.

File fsync is the durability boundary; this version does not fsync directory entries, so creation/rotation names are not promised durable through power loss. The legacy NewDLServer constructor remains as a panic-on-open-error wrapper; new code should use Open. Log, Flush, NewFile and Close now return errors which callers should check.

Run `go run .` for the demo and `go test -race ./...` for buffered retrieval, rotation, restart, concurrent writing, recovery and failure-state tests.

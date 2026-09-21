package durablelogs

import (
	"bufio"
	"durablelogs/durablelogs/pb"
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"os"
	"sync"
	"time"
)

var ErrClosed = errors.New("durablelogs: closed")
var ErrPoisoned = errors.New("durablelogs: poisoned")

type DurableLogger struct {
	directory                         string
	maxPerFile                        int
	bufWriter                         *bufio.Writer
	currentFile                       *os.File
	mu                                sync.Mutex
	currentFileNum, logsInCurrentFile int
	pending                           []string
	closed                            bool
	poison                            error
}

func Open(directory string, maxPerFile int) (*DurableLogger, error) {
	if maxPerFile < 1 {
		return nil, errors.New("maxPerFile must be positive")
	}
	if e := os.MkdirAll(directory, 0700); e != nil {
		return nil, e
	}
	ids, e := segments(directory)
	if e != nil {
		return nil, e
	}
	last, count := 0, 0
	for i, id := range ids {
		records, e := readSegment(segmentPath(directory, id), i == len(ids)-1)
		if e != nil {
			return nil, fmt.Errorf("segment %d: %w", id, e)
		}
		last = id
		count = len(records)
	}
	f, e := os.OpenFile(segmentPath(directory, last), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if e != nil {
		return nil, e
	}
	return &DurableLogger{directory: directory, maxPerFile: maxPerFile, currentFile: f, bufWriter: bufio.NewWriter(f), currentFileNum: last, logsInCurrentFile: count}, nil
}

// NewDLServer is the legacy panic-on-open-error wrapper. New code should use Open.
func NewDLServer(directory string, maxPerFile int) *DurableLogger {
	d, e := Open(directory, maxPerFile)
	if e != nil {
		panic(e)
	}
	return d
}
func (d *DurableLogger) check() error {
	if d.closed {
		return ErrClosed
	}
	if d.poison != nil {
		return errors.Join(ErrPoisoned, d.poison)
	}
	return nil
}
func (d *DurableLogger) fail(e error) error {
	if e != nil {
		d.poison = e
	}
	return e
}
func (d *DurableLogger) Log(message string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e := d.check(); e != nil {
		return e
	}
	if len(message) > maxRecord-1024 {
		return errors.New("message too large")
	}
	b, e := proto.Marshal(&pb.Log{Log: message, Timestamp: time.Now().UTC().Format(time.RFC3339Nano)})
	if e != nil {
		return e
	}
	if d.logsInCurrentFile >= d.maxPerFile {
		if e = d.rotate(); e != nil {
			return d.fail(e)
		}
	}
	var h [4]byte
	binary.LittleEndian.PutUint32(h[:], uint32(len(b)))
	if _, e = d.bufWriter.Write(h[:]); e == nil {
		_, e = d.bufWriter.Write(b)
	}
	if e != nil {
		return d.fail(e)
	}
	d.pending = append(d.pending, message)
	d.logsInCurrentFile++
	return nil
}
func (d *DurableLogger) flush() error {
	if e := d.bufWriter.Flush(); e != nil {
		return d.fail(e)
	}
	if e := d.currentFile.Sync(); e != nil {
		return d.fail(e)
	}
	d.pending = nil
	return nil
}
func (d *DurableLogger) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e := d.check(); e != nil {
		return e
	}
	return d.flush()
}
func (d *DurableLogger) rotate() error {
	if e := d.flush(); e != nil {
		return e
	}
	next := d.currentFileNum + 1
	f, e := os.OpenFile(segmentPath(d.directory, next), os.O_CREATE|os.O_EXCL|os.O_RDWR|os.O_APPEND, 0600)
	if e != nil {
		return e
	}
	if e = d.currentFile.Close(); e != nil {
		f.Close()
		return e
	}
	d.currentFile = f
	d.bufWriter = bufio.NewWriter(f)
	d.currentFileNum = next
	d.logsInCurrentFile = 0
	return nil
}
func (d *DurableLogger) NewFile() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e := d.check(); e != nil {
		return e
	}
	return d.fail(d.rotate())
}
func (d *DurableLogger) GetBufferedLogs() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.pending...)
}
func (d *DurableLogger) GetCurrentFile() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.currentFileNum
}
func (d *DurableLogger) ReadAll() ([]*pb.Log, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e := d.check(); e != nil {
		return nil, e
	}
	if e := d.flush(); e != nil {
		return nil, e
	}
	ids, e := segments(d.directory)
	if e != nil {
		return nil, e
	}
	var result []*pb.Log
	for _, id := range ids {
		records, e := readSegment(segmentPath(d.directory, id), false)
		if e != nil {
			return nil, e
		}
		result = append(result, records...)
	}
	return result, nil
}
func (d *DurableLogger) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return nil
	}
	d.closed = true
	return errors.Join(d.poison, d.flush(), d.currentFile.Close())
}

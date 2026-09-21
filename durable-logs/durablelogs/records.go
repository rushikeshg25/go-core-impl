package durablelogs

import (
	"durablelogs/durablelogs/pb"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const maxRecord = 64 << 20

func segments(dir string) ([]int, error) {
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	ids := []int{}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "dl-") {
			continue
		}
		n, e := strconv.Atoi(strings.TrimPrefix(entry.Name(), "dl-"))
		if e != nil || n < 0 || entry.IsDir() || entry.Name() != fmt.Sprintf("dl-%d", n) {
			return nil, fmt.Errorf("invalid segment %q", entry.Name())
		}
		ids = append(ids, n)
	}
	sort.Ints(ids)
	return ids, nil
}
func segmentPath(dir string, n int) string { return filepath.Join(dir, fmt.Sprintf("dl-%d", n)) }
func readSegment(path string, repair bool) ([]*pb.Log, error) {
	flags := os.O_RDONLY
	if repair {
		flags = os.O_RDWR
	}
	f, e := os.OpenFile(path, flags, 0600)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		return nil, e
	}
	end := stat.Size()
	var records []*pb.Log
	off := int64(0)
	for off < end {
		valid := off
		var h [4]byte
		_, e = f.ReadAt(h[:], off)
		n := int64(binary.LittleEndian.Uint32(h[:]))
		if e == nil && n > maxRecord {
			return nil, fmt.Errorf("oversized record at %d", off)
		}
		if e != nil || end-off-4 < n {
			if !repair {
				return nil, io.ErrUnexpectedEOF
			}
			if e = f.Truncate(valid); e != nil {
				return nil, e
			}
			return records, f.Sync()
		}
		b := make([]byte, n)
		if n > 0 {
			if _, e = f.ReadAt(b, off+4); e != nil {
				return nil, e
			}
		}
		record := new(pb.Log)
		if e = proto.Unmarshal(b, record); e != nil {
			return nil, e
		}
		records = append(records, record)
		off += 4 + n
	}
	return records, nil
}

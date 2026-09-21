package durablelogs

import (
	"errors"
	"os"
	"sync"
	"testing"
)

func TestRotateReopenAndBuffer(t *testing.T) {
	dir := t.TempDir()
	d, e := Open(dir, 2)
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []string{"a", "b", "c", "d", "e"} {
		if e = d.Log(s); e != nil {
			t.Fatal(e)
		}
	}
	pending := d.GetBufferedLogs()
	if len(pending) != 1 || pending[0] != "e" {
		t.Fatal(pending)
	}
	pending[0] = "changed"
	if d.GetBufferedLogs()[0] != "e" {
		t.Fatal("aliased buffer")
	}
	if e = d.Close(); e != nil {
		t.Fatal(e)
	}
	d, e = Open(dir, 2)
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	records, e := d.ReadAll()
	if e != nil || len(records) != 5 || records[4].Log != "e" {
		t.Fatal(records, e)
	}
	if d.GetCurrentFile() != 2 {
		t.Fatal("wrong latest segment")
	}
}
func TestConcurrentLog(t *testing.T) {
	d, _ := Open(t.TempDir(), 7)
	defer d.Close()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				if e := d.Log("message"); e != nil {
					t.Error(e)
				}
			}
		}()
	}
	wg.Wait()
	records, e := d.ReadAll()
	if e != nil || len(records) != 200 {
		t.Fatal(len(records), e)
	}
}
func TestTailRecoveryAndClosed(t *testing.T) {
	dir := t.TempDir()
	d, _ := Open(dir, 2)
	d.Log("good")
	d.Close()
	f, _ := os.OpenFile(segmentPath(dir, 0), os.O_APPEND|os.O_WRONLY, 0600)
	f.Write([]byte{1, 2})
	f.Close()
	d, e := Open(dir, 2)
	if e != nil {
		t.Fatal(e)
	}
	records, e := d.ReadAll()
	if e != nil || len(records) != 1 {
		t.Fatal(records, e)
	}
	d.Close()
	if d.Close() != nil || !errors.Is(d.Log("bad"), ErrClosed) {
		t.Fatal("close lifecycle")
	}
}
func TestPoisonOnFlush(t *testing.T) {
	d, _ := Open(t.TempDir(), 2)
	d.Log("a")
	d.currentFile.Close()
	if d.Flush() == nil {
		t.Fatal("ignored disk error")
	}
	if !errors.Is(d.Log("b"), ErrPoisoned) {
		t.Fatal("not poisoned")
	}
	if d.Close() == nil {
		t.Fatal("lost close error")
	}
}

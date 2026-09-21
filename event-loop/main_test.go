package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDrainAndBound(t *testing.T) {
	el := NewEventLoop()
	var active, peak, tasks, callbacks atomic.Int32
	for i := 0; i < 100; i++ {
		if !el.AddEvent(&Event{Async: true, Task: func() {
			n := active.Add(1)
			for p := peak.Load(); n > p; p = peak.Load() {
				if peak.CompareAndSwap(p, n) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			tasks.Add(1)
			active.Add(-1)
		}, Callback: func() { callbacks.Add(1) }}) {
			t.Fatal("rejected")
		}
	}
	el.StopEventLoop()
	el.StopEventLoop()
	if tasks.Load() != 100 || callbacks.Load() != 100 || peak.Load() > 5 {
		t.Fatalf("tasks %d callbacks %d peak %d", tasks.Load(), callbacks.Load(), peak.Load())
	}
	if el.AddEvent(&Event{}) {
		t.Fatal("accepted after stop")
	}
}
func TestConcurrentSubmitStop(t *testing.T) {
	el := NewEventLoop()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); el.AddEvent(&Event{Task: func() {}}) }()
	}
	el.StopEventLoop()
	wg.Wait()
}
func TestSyncCallbackOrder(t *testing.T) {
	el := NewEventLoop()
	var n int
	el.AddEvent(&Event{Task: func() { n = 1 }, Callback: func() {
		if n != 1 {
			t.Error("callback before task")
		}
		n = 2
	}})
	el.StopEventLoop()
	if n != 2 {
		t.Fatal(n)
	}
}

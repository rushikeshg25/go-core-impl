package main

import (
	"context"
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

func TestCancellationStillCompletesCallbacks(t *testing.T) {
	for _, async := range []bool{false, true} {
		el := NewEventLoop()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var tasks, callbacks atomic.Int32
		if !el.AddEvent(&Event{Async: async, Context: ctx, Task: func() { tasks.Add(1) }, Callback: func() { callbacks.Add(1) }}) {
			t.Fatal("rejected")
		}
		if !AddCallback(el, &Event{Context: ctx, TaskContext: func(context.Context) { tasks.Add(1) }, Callback: func() { callbacks.Add(1) }}) {
			t.Fatal("rejected callback")
		}
		el.StopEventLoop()
		if tasks.Load() != 0 || callbacks.Load() != 2 {
			t.Fatalf("tasks %d callbacks %d", tasks.Load(), callbacks.Load())
		}
	}
}

func TestRunningTaskCooperatesWithCancellation(t *testing.T) {
	el := NewEventLoop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	var callbacks atomic.Int32
	el.AddEvent(&Event{Async: true, Context: ctx, TaskContext: func(taskContext context.Context) {
		close(started)
		<-taskContext.Done()
	}, Callback: func() { callbacks.Add(1) }})
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("task did not start")
	}
	cancel()
	done := make(chan struct{})
	go func() { el.StopEventLoop(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("canceled task did not drain")
	}
	if callbacks.Load() != 1 {
		t.Fatalf("callbacks %d", callbacks.Load())
	}
}

func TestInvalidAndStoppedCallbackAdmission(t *testing.T) {
	el := NewEventLoop()
	if el.AddEvent(nil) || AddCallback(el, nil) {
		t.Fatal("accepted nil event")
	}
	el.StopEventLoop()
	if AddCallback(el, &Event{}) {
		t.Fatal("accepted callback after stop")
	}
}

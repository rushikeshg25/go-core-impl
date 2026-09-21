package main

import (
	"context"
	"fmt"
	"sync"
)

// Event runs Task and then Callback. Async tasks use a bounded pool.
type Event struct {
	// Context defaults to Background. Cancellation before execution skips the task.
	Context context.Context
	// TaskContext takes precedence over Task and cooperates with cancellation.
	TaskContext func(context.Context)
	Task        func()
	Callback    func()
	Async       bool
	isAsync     bool
}
type EventLoop struct {
	Events    chan Event
	Callbacks chan Event
	once      sync.Once
	mu        sync.Mutex
	stopped   bool
	stop      chan struct{}
	done      chan struct{}
	wg        sync.WaitGroup
}

func NewEventLoop() *EventLoop {
	return &EventLoop{Events: make(chan Event, 5), Callbacks: make(chan Event, 5), stop: make(chan struct{}), done: make(chan struct{})}
}

// Start is idempotent and returns immediately. The wait group completes after StopEventLoop.
func (el *EventLoop) Start() *sync.WaitGroup {
	el.once.Do(func() { el.wg.Add(1); go el.run() })
	return &el.wg
}
func invoke(f func()) {
	if f != nil {
		f()
	}
}
func runTask(e Event) {
	ctx := e.Context
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return
	}
	if e.TaskContext != nil {
		e.TaskContext(ctx)
	} else {
		invoke(e.Task)
	}
}

func (el *EventLoop) run() {
	defer el.wg.Done()
	defer close(el.done)
	completed := make(chan Event, 5)
	active := 0
	stopping := false
	stop := el.stop
	for {
		if stopping && active == 0 && len(el.Events) == 0 && len(el.Callbacks) == 0 {
			return
		}
		events := el.Events
		if active == 5 {
			events = nil
		}
		select {
		case e := <-events:
			if e.Async || e.isAsync {
				active++
				go func() { runTask(e); completed <- e }()
			} else {
				runTask(e)
				invoke(e.Callback)
			}
		case e := <-completed:
			active--
			invoke(e.Callback)
		case e := <-el.Callbacks:
			runTask(e)
			invoke(e.Callback)
		case <-stop:
			stopping = true
			stop = nil
		}
	}
}

// AddEvent returns false after shutdown; accepted events are drained on shutdown.
func (el *EventLoop) AddEvent(event *Event) bool {
	if event == nil {
		return false
	}
	el.Start()
	el.mu.Lock()
	defer el.mu.Unlock()
	if el.stopped {
		return false
	}
	el.Events <- *event
	return true
}

// StopEventLoop drains accepted work. Call from outside tasks and callbacks.
func (el *EventLoop) StopEventLoop() {
	el.Start()
	el.mu.Lock()
	if !el.stopped {
		el.stopped = true
		close(el.stop)
	}
	el.mu.Unlock()
	<-el.done
}
func AddCallback(el *EventLoop, event *Event) bool {
	if event == nil {
		return false
	}
	el.Start()
	el.mu.Lock()
	defer el.mu.Unlock()
	if el.stopped {
		return false
	}
	el.Callbacks <- *event
	return true
}
func main() {
	el := NewEventLoop()
	el.AddEvent(&Event{Async: true, Task: func() { fmt.Println("task") }, Callback: func() { fmt.Println("done") }})
	el.StopEventLoop()
}

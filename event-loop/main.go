package main

import (
	"sync"
)

type Event struct {
	Task     func()
	Callback func()
	isAsync  bool
}

type EventLoop struct {
	Events    chan Event
	Callbacks chan Event
	stop      chan bool
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		Events:    make(chan Event, 5),
		Callbacks: make(chan Event, 5),
		stop:      make(chan bool),
	}
}

func (el *EventLoop) Start() {
	var wg sync.WaitGroup
	pool := make(chan struct{}, 5) //For async tasks

	wg.Add(1)
	for {
		defer wg.Done()
		select {
		case e := <-el.Events:
			if e.isAsync {
				pool <- struct{}{}
				go func() {
					defer func() {
						<-pool
					}()
					e.Task()

					if e.Callback != nil {
						AddCallback(
							el,
							&Event{
								Task: e.Callback,
							},
						)
					}
				}()
			}
		}
	}

}

func AddEvent(el *EventLoop, event *Event) {
	el.Events <- *event
}

func AddCallback(el *EventLoop, event *Event) {
	el.Callbacks <- *event
}

func StopEventLoop(el *EventLoop) {
	el.stop <- true
}

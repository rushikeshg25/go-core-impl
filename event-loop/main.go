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

func (el *EventLoop) Start() *sync.WaitGroup {
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
			} else {
				e.Task()
			}
		case e := <-el.Callbacks:
			e.Task()
		case stop := <-el.stop:
			if stop {
				return
			}
		}
	}
	return &wg
}

func (el *EventLoop) AddEvent(event *Event) {
	el.Events <- *event
}

func (el *EventLoop) StopEventLoop() {
	el.stop <- true
}

func AddCallback(el *EventLoop, event *Event) {
	el.Callbacks <- *event
}

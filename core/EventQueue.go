package core

import (
	"github.com/gw123/GMQ/core/interfaces"
)

const MaxEventLen = 4096

type EventQueue struct {
	Events chan interfaces.Event
	len    int
	app    interfaces.App
}

func NewEventQueue(app interfaces.App) *EventQueue {
	this := new(EventQueue)
	this.len = MaxEventLen
	this.Events = make(chan interfaces.Event, MaxEventLen)
	this.app = app
	return this
}

func (q *EventQueue) Push(event interfaces.Event) error {
	q.Events <- event
	return nil
}

func (q *EventQueue) Pop() (interfaces.Event, error) {
	return <-q.Events, nil
}

package core

import (
	"github.com/gw123/GMQ/core/interfaces"
)

const MaxEventLen = 4096

type EventQueue struct {
	Events chan interfaces.Msg
	len    int
	app    interfaces.App
}

func NewEventQueue(app interfaces.App) *EventQueue {
	this := new(EventQueue)
	this.len = MaxEventLen
	this.Events = make(chan interfaces.Msg, MaxEventLen)
	this.app = app
	return this
}

func (q *EventQueue) Push(event interfaces.Msg) error {
	q.Events <- event
	return nil
}

func (q *EventQueue) Pop() (interfaces.Msg, error) {
	return <-q.Events, nil
}

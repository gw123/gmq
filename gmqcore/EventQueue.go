package gmqcore

import (
	"github.com/gw123/gmq"
)

const MaxEventLen = 4096

type EventQueue struct {
	Events chan gmq.Msg
	len    int
	app    gmq.App
}

func NewEventQueue(app gmq.App) *EventQueue {
	this := new(EventQueue)
	this.len = MaxEventLen
	this.Events = make(chan gmq.Msg, MaxEventLen)
	this.app = app
	return this
}

func (q *EventQueue) Push(event gmq.Msg) error {
	q.Events <- event
	return nil
}

func (q *EventQueue) Pop() (gmq.Msg, error) {
	return <-q.Events, nil
}

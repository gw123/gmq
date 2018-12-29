package app

import (
	"sync"
	"errors"
	"github.com/gw123/GMQ/interfaces"
)

const MaxEventLen = 1024

type EventQueue struct {
	Events []interfaces.Event
	head   int
	tail   int
	len    int
	Mutex  sync.Mutex
	app    interfaces.App
}

func NewEventQueue(app interfaces.App) *EventQueue {
	this := new(EventQueue)
	this.head = 0
	this.tail = 0
	this.len = 0
	this.Events = make([]interfaces.Event, MaxEventLen)
	this.app = app
	return this
}

func (this *EventQueue) Push(event interfaces.Event) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.len == MaxEventLen {
		return errors.New("队列已满")
	}

	if this.head >= MaxEventLen {
		this.head = 0
	}
	this.Events[this.head] = event
	this.head++;
	this.head = this.head % MaxEventLen
	this.len++;
	//this.app.Debug("EventQueue", "Push "+event.GetMsgId()+" : "+event.GetEventName())
	return nil
}

func (this *EventQueue) Pop() (interfaces.Event, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	//this.app.Debug("EventQueue","eee")
	//this.app.Debug("EventQueue", fmt.Sprintf("%s",this.app))
	if this.len == 0 {
		return nil, errors.New("队列为空")
	}

	event := this.Events[this.tail]
	this.Events[this.tail] = nil
	this.tail++
	this.len--;
	this.tail = this.tail % MaxEventLen
	this.app.Debug("EventQueue", "Pop "+event.GetMsgId()+" : "+event.GetEventName())
	return event, nil
}

package common

import (
	"time"
	"fmt"
)

type Event struct {
	MsgId      string
	EventName  string
	ModuleName string
	Payload    string
}

func NewEvent(eventType string, result string) *Event {
	this := new(Event)
	this.MsgId = this.CreateMsgId()
	this.Payload = result
	this.EventName = eventType
	return this
}

func (this *Event) GetMsgId() string {
	return this.MsgId
}

func (this *Event) GetEventName() string {
	return this.EventName
}

func (this *Event) GetPayload() []byte {
	return []byte(this.Payload)
}

func (this *Event) CreateMsgId() string {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	nano := now.UnixNano() % 1000000
	date := fmt.Sprintf("%d:%d:%d:%d", hour, minute, second, nano)
	return fmt.Sprintf("%s", date)
}

func (this *Event) SetSourceModule(name string) {
	this.ModuleName = name
}

func (this *Event) GetSourceModule() string {
	return this.ModuleName
}

type RequestEvent struct {
	Event
	Token string
}

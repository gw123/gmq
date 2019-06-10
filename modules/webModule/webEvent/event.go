package webEvent

import (
	"time"
	"fmt"
)

type Event struct {
	MsgId           string
	EventName       string
	SourcModuleName string
	DstModuleName   string
	Payload         string
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

func (this *Event) GetSourceModule() string {
	return this.SourcModuleName
}

func (this *Event) GetDstModule() string {
	return this.DstModuleName
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

type RequestEvent struct {
	Event
	Token string
}

package common_types

import (
	"fmt"
	"time"
)

type Event struct {
	MsgId            string
	EventName        string
	Payload          []byte
	sourceModuleName string
	dstModuleName    string
}

func NewEvent(eventType string, result []byte) *Event {
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
	return this.Payload
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
	this.sourceModuleName = name
}

func (this *Event) GetSourceModule() string {
	return this.sourceModuleName
}

func (this *Event) GetDstModule() string {
	if this.dstModuleName == "" {
		return "*"
	}
	return this.dstModuleName
}

func (this *Event) SetDstModule(string2 string) {
	this.dstModuleName = string2
}

package common_types

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Event struct {
	MsgId            string
	EventName        string
	Payload          interface{}
	sourceModuleName string
	dstModuleName    string
}

func NewEvent(eventType string, payload interface{}) *Event {
	this := new(Event)
	this.MsgId = this.CreateMsgId()
	this.Payload = payload
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
	ret, err := json.Marshal(this.Payload)
	if err != nil {
		return []byte{}
	}
	return ret
}

func (this *Event) GetInterface() interface{} {
	return this.Payload
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

func (this *Event) SetSourceModule(name string) {
	this.sourceModuleName = name
}
func (this *Event) SetDstModule(string2 string) {
	this.dstModuleName = string2
}

func (this *Event) CreateMsgId() string {
	now := time.Now()
	rand := rand.Int31n(100000)
	minute := now.Minute()
	second := now.Second()
	nano := now.UnixNano() % 1000000
	date := fmt.Sprintf("%d:%d:%d:%d", rand, minute, second, nano)
	return fmt.Sprintf("%s", date)
}
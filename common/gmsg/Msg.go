package gmsg

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Msg struct {
	MsgId            string
	EventName        string
	Payload          interface{}
	sourceModuleName string
	dstModuleName    string
}

func NewEvent(eventType string, payload interface{}) *Msg {
	this := new(Msg)
	this.MsgId = this.CreateMsgId()
	this.Payload = payload
	this.EventName = eventType
	return this
}

func (this *Msg) GetMsgId() string {
	return this.MsgId
}

func (this *Msg) GetEventName() string {
	return this.EventName
}

func (this *Msg) GetPayload() []byte {
	ret, err := json.Marshal(this.Payload)
	if err != nil {
		return []byte{}
	}
	return ret
}

func (this *Msg) GetInterface() interface{} {
	return this.Payload
}

func (this *Msg) GetSourceModule() string {
	return this.sourceModuleName
}

func (this *Msg) GetDstModule() string {
	if this.dstModuleName == "" {
		return "*"
	}
	return this.dstModuleName
}

func (this *Msg) SetSourceModule(name string) {
	this.sourceModuleName = name
}

func (this *Msg) SetDstModule(string2 string) {
	this.dstModuleName = string2
}

func (this *Msg) CreateMsgId() string {
	now := time.Now()
	rand := rand.Int31n(100000)
	minute := now.Minute()
	second := now.Second()
	nano := now.UnixNano() % 1000000
	date := fmt.Sprintf("%d:%d:%d:%d", rand, minute, second, nano)
	return fmt.Sprintf("%s", date)
}
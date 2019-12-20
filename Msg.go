package gmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Msg interface {
	GetMsgId() string
	GetEventName() string
	GetPayload() []byte
	GetInterface() interface{}
	GetSourceModule() string
	GetDstModule() string
	GetContext() context.Context
}

type BaseMsg struct {
	MsgId            string
	EventName        string
	Payload          interface{}
	sourceModuleName string
	dstModuleName    string
}

func NewEvent(eventType string, payload interface{}) *BaseMsg {
	this := new(BaseMsg)
	this.MsgId = this.CreateMsgId()
	this.Payload = payload
	this.EventName = eventType
	return this
}

func (this *BaseMsg) GetMsgId() string {
	return this.MsgId
}

func (this *BaseMsg) GetEventName() string {
	return this.EventName
}

func (this *BaseMsg) GetPayload() []byte {
	ret, err := json.Marshal(this.Payload)
	if err != nil {
		return []byte{}
	}
	return ret
}

func (this *BaseMsg) GetInterface() interface{} {
	return this.Payload
}

func (this *BaseMsg) GetSourceModule() string {
	return this.sourceModuleName
}

func (this *BaseMsg) GetDstModule() string {
	if this.dstModuleName == "" {
		return "*"
	}
	return this.dstModuleName
}

func (this *BaseMsg) SetSourceModule(name string) {
	this.sourceModuleName = name
}

func (this *BaseMsg) SetDstModule(string2 string) {
	this.dstModuleName = string2
}

func (this *BaseMsg) CreateMsgId() string {
	now := time.Now()
	rand := rand.Int31n(100000)
	minute := now.Minute()
	second := now.Second()
	nano := now.UnixNano() % 1000000
	date := fmt.Sprintf("%d:%d:%d:%d", rand, minute, second, nano)
	return fmt.Sprintf("%s", date)
}

type UpdateCacheMsg struct {
	BaseMsg
	Cachekey  CacheKey
	Arguments []interface{}
}

func NewUpdateCacheMsg(cachekey CacheKey, argument []interface{}) *UpdateCacheMsg {
	this := new(UpdateCacheMsg)
	this.EventName = "updateCache"
	this.Cachekey = cachekey
	this.Arguments = argument
	this.MsgId = this.GetMsgId()
	return this
}

func (l *UpdateCacheMsg) GetMsgId() string {
	return l.MsgId
}

func (l *UpdateCacheMsg) GetEventName() string {
	return l.EventName
}

func (l *UpdateCacheMsg) GetPayload() []byte {
	return l.Payload.([]byte)
}

func (l *UpdateCacheMsg) GetInterface() interface{} {
	return l.Payload
}

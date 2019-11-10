package gmsg

import "github.com/gw123/GMQ/core/interfaces"

type UpdateCacheMsg struct {
	Msg
	Cachekey  interfaces.CacheKey
	Arguments []interface{}
}

func NewUpdateCacheMsg(cachekey interfaces.CacheKey, argument []interface{}) *UpdateCacheMsg {
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

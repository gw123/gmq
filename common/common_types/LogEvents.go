package common_types

type LogEvent struct {
	Event
}

func NewLogEvent(result []byte) *LogEvent {
	this := new(LogEvent)
	this.Payload = result
	this.EventName = "log"
	this.MsgId = this.GetMsgId()
	return this
}

func (this *LogEvent) GetMsgId() string {
	return this.MsgId
}

func (this *LogEvent) GetEventName() string {
	return this.EventName
}

func (this *LogEvent) GetPayload() []byte {
	return this.Payload.([]byte)
}

func (this *LogEvent) GetInterface() interface{} {
	return this.Payload
}

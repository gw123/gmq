package gmsg

type LogMsg struct {
	Msg
}

func NewLogMsg(result []byte) *LogMsg {
	this := new(LogMsg)
	this.Payload = result
	this.EventName = "log"
	this.MsgId = this.GetMsgId()
	return this
}

func (l *LogMsg) GetMsgId() string {
	return l.MsgId
}

func (l *LogMsg) GetEventName() string {
	return l.EventName
}

func (l *LogMsg) GetPayload() []byte {
	return l.Payload.([]byte)
}

func (l *LogMsg) GetInterface() interface{} {
	return l.Payload
}

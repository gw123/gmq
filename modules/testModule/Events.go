package testModule

type ResultEvent struct {
	MsgId     string
	EventType string
	Payload   string
}

func NewResultEvent(msgId ,result string) *ResultEvent {
	this := new(ResultEvent)
	this.MsgId = msgId
	this.Payload = result
	this.EventType = "Test"
	return this
}

func (this *ResultEvent) GetMsgId() string {
	return this.MsgId
}

func (this *ResultEvent) GetEventName() string {
	return this.EventType
}

func (this *ResultEvent) GetPayload() []byte {
	return []byte(this.Payload)
}

package app

type Event struct {
	MsgId     string
	EventType string
	Payload   string
}

func (this *Event) GetMsgId() string {
	return this.MsgId
}

func (this *Event) GetEventType() string {
	return this.EventType
}

func (this *Event) GetPayload() string {
	return this.Payload
}

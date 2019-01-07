package types

//上报阿里云任务执行结果
//payload 是 Lhmsg 的字符串
type ResultEvent struct {
	Event
	Payload      string
}

func NewResultEvent(result []byte) *ResultEvent {
	this := new(ResultEvent)
	this.MsgId = this.GetMsgId()
	this.Payload = string(result)
	this.EventName = "reply"
	return this
}

func (this *ResultEvent) GetMsgId() string {
	return this.MsgId
}

func (this *ResultEvent) GetEventName() string {
	return this.EventName
}

func (this *ResultEvent) GetPayload() []byte {
	return []byte(this.Payload)
}

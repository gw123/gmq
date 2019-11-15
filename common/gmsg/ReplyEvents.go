package gmsg

//上报阿里云任务执行结果
//payload 是 Lhmsg 的字符串
type ResultEvent struct {
	Msg
	Payload string
}

func NewResultEvent(result []byte) *ResultEvent {
	this := new(ResultEvent)
	this.MsgId = this.GetMsgId()
	this.Payload = string(result)
	this.EventName = "reply"
	return this
}

func (result *ResultEvent) GetMsgId() string {
	return result.MsgId
}

func (result *ResultEvent) GetEventName() string {
	return result.EventName
}

func (result *ResultEvent) GetPayload() []byte {
	return []byte(result.Payload)
}

func (result *ResultEvent) GetInterface() interface{} {
	return result.Payload
}

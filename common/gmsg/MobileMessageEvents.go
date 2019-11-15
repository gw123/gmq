package gmsg

type MobileMessageEvent struct {
	Msg
	Code    string
	Modbile string
}

func NewMobileMessageEvent(code, mobile string) *MobileMessageEvent {
	this := new(MobileMessageEvent)
	this.Code = code
	this.Modbile = mobile
	this.EventName = "sendMobileMessage"
	this.MsgId = this.GetMsgId()
	return this
}

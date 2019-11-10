package interfaces

type Msg interface {
	GetMsgId() string
	GetEventName() string
	GetPayload() []byte
	GetInterface() interface{}
	GetSourceModule() string
	GetDstModule() string
}

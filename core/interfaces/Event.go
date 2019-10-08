package interfaces

type Event interface {
	GetMsgId() string
	GetEventName() string
	GetPayload() []byte
	GetInterface() interface{}
	GetSourceModule() string
	GetDstModule() string
}

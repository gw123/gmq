package interfaces

type Event interface {
	GetMsgId() string
	GetEventName() string
	GetPayload() []byte
	GetSourceModule() string
	GetDstModule() string
}

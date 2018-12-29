package interfaces

type Event interface {
	GetMsgId() string
	GetEventName() string
	GetPayload() []byte
	SetSourceModule(string)
	GetSourceModule() string
}

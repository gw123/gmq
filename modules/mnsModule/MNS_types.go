package mnsModule

type MNSmsg struct {
	Payload     string    `json:"payload"`
	MessageType string    `json:"messagetype"`
	Topic       string    `json:"topic"`
	Messageid   uint64    `json:"messageid"`
	Timestamp   uint32 `json:"timestamp"`
}

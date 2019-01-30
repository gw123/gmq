package mnsModule

import "time"

type MNSmsg struct {
	Payload     string    `json:"payload"`
	MessageType string    `json:"messagetype"`
	Topic       string    `json:"topic"`
	Messageid   uint64    `json:"messageid"`
	Timestamp   time.Time `json:"timestamp"`
}

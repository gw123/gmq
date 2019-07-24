package models

import (
	"time"
)

type PingLog struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	ClientSendAt time.Time
	Latency      uint
	Ip           string
	City         string
	Ttl          uint
	ClientId     string
	BytesIn      uint
	BytesOut     uint
	Payload      string    `json:"payload" form:"payload" gorm:"size:1024"`
}

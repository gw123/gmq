package db_models

import (
	"time"
)

type Client struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	Name      string     `gorm:"size:32" json:"name"`
	Info      string     `gorm:"size:1024" json:"info"`
	Token     string     `json:"token"`
	Secret    string     `json:"secret"`
	ClientTasks []ClientTask `json:"client_tasks"`
}

func (u Client) TableName() string {
	return "clients"
}

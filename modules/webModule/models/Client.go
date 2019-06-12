package models

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	Name    string `gorm:"size:32"`
	Info    string `gorm:"size:1024"`
	Token   string
	Secret   string
}

func (u Client) TableName() string {
	return "clients"
}

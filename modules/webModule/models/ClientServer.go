package models

import "github.com/jinzhu/gorm"

type ClientServer struct {
	gorm.Model
	ServerId      uint
	ServerName    string
	ServerVersion string
	Status        uint8
}

func (u ClientServer) TableName() string {
	return "client_servers"
}

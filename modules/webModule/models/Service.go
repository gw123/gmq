package models

import "github.com/jinzhu/gorm"

type Service struct {
	gorm.Model
	Name    string `json:"name" form:"name" gorm:"size:32" validate:"required"`
	Desc    string `json:"desc" form:"desc" gorm:"size:255" validate:"required"`
	Version string `json:"version" form:"version" gorm:"size:32" validate:"required"`
	Config  string `json:"config" form:"config" gorm:"size:2048" validate:"required"`
	File    string `json:"file" form:"file" gorm:"size:128" validate:"required"`
}

func (u Service) TableName() string {
	return "services"
}

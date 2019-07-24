package models

import "github.com/jinzhu/gorm"

type TaskDetail struct {
	gorm.Model
	TaskId   string `json:"task_id"  form:"task_id"`
	Name     string `json:"name" form:"name" gorm:"size:32" validate:"required"`
	Version  string `json:"version" form:"version" gorm:"size:32" validate:"required"`
	Config   string `json:"config" form:"config" gorm:"size:2048" `
	File     string `json:"file" form:"file" gorm:"size:128" validate:"required"`
	FileType string `json:"file_type" form:"file_type" gorm:"size:128" validate:"required"`
	CheckSum string `json:"checksum" form:"checksum" gorm:"size:128" `
}

func (u TaskDetail) TableName() string {
	return "task_details"
}

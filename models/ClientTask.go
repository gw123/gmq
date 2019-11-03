package models

import "github.com/jinzhu/gorm"

type ClientTask struct {
	gorm.Model
	ClientId    int    `json:"client_id" form:"client_id" validate:"required"`
	TaskId      int    `json:"task_id" form:"task_id" validate:"required"`
	TaskName    string `json:"task_name" form:"task_name" `
	TaskVersion string `json:"task_version" form:"task_version" validate:"required"`
	Status      uint8  `json:"status" form:"status"`
}

func (u ClientTask) TableName() string {
	return "client_tasks"
}

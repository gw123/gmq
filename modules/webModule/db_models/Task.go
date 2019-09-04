package db_models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Name     string `json:"name" form:"name" gorm:"size:32" validate:"required"`
	Desc     string `json:"desc" form:"desc" gorm:"size:255"`
	Type     string `json:"type" form:"type" gorm:"size:255"`
	//Version  string `json:"version" form:"version" gorm:"size:32" validate:"required"`
	//Config   string `json:"config" form:"config" gorm:"size:2048" `
	//File     string `json:"file" form:"file" gorm:"size:128" validate:"required"`
	//FileType string `json:"file_type" form:"file_type" gorm:"size:128" validate:"required"`
	//CheckSum string `json:"checksum" form:"checksum" gorm:"size:128" `
}

func (u *Task) TableName() string {
	return "tasks"
}

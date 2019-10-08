package models

import "time"

type Resource struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Title  string `json:"title"`
	UserId string `json:"user_id"`
	Covers string `json:"covers"`
	Type   string `json:"type"`
}

func (u Resource) TableName() string {
	return "resource"
}

type Article struct {
	Rid       int       `json:"rid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (a Article) TableName() string {
	return "type_article"
}

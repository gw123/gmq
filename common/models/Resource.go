package models

import (
	"encoding/json"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Category struct {
	ID    int    `gorm:"primary_key" json:"id"`
	Title string `json:"title"`
}

func (Category) TableName() string {
	return "category"
}

type Tag struct {
	ID    int    `gorm:"primary_key" json:"id"`
	Title string `json:"title"`
}

func (Tag) TableName() string {
	return "tag"
}

type TagGroup struct {
	Model
	TagId   uint `json:"tag_id"`
	GroupId uint `json:"group_id"`
}

func (TagGroup) TableName() string {
	return "tag_group"
}

type Group struct {
	Model
	Title      string    `gorm:"column:title" json:"title"`
	UserId     uint      `gorm:"column:user_id"json:"user_id"`
	Covers     []string  `gorm:"-" json:"covers"`
	CoversStr  string    `gorm:"column:covers" json:"covers_str"`
	Desc       string    `gorm:"column:desc" json:"desc"`
	CategoryId uint      `gorm:"column:category_id" json:"category_id"`
	TagIds     []uint    `gorm:"-" json:"tags"`
	Display    string    `json:"display"`
	Chapters   []Chapter `json:"chapters"`
}

func (Group) TableName() string {
	return "group"
}

func (g *Group) BeforeSave() (err error) {
	var str []byte
	if str, err = json.Marshal(g.Covers); err != nil {
		return err
	}
	g.CoversStr = string(str)
	return
}

type Chapter struct {
	Model
	Title     string     `json:"title"`
	GroupId   uint        `json:"group_id"`
	ParentId  uint        `json:"parent_id"`
	Resources []Resource `json:"resources"`
}

func (Chapter) TableName() string {
	return "chapter"
}

type Resource struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Title     string     `json:"title"`
	UserId    uint       `json:"user_id"`
	GroupId   uint       `json:"group_id"`
	ChapterId uint       `json:"chapter_id"`
	Covers    string     `json:"covers"`
	Type      string     `json:"type"`
	Article   *Article   `json:"article" gorm:"ForeignKey:rid"`
	Testpaper *Testpaper `json:"testpaper" gorm:"ForeignKey:rid"`
}

func (r Resource) TableName() string {
	return "resource"
}

type Article struct {
	Id        uint      `json:"id"`
	Rid       int       `json:"rid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (a Article) TableName() string {
	return "type_article"
}

type Testpaper struct {
	Id        uint      `json:"id"`
	Rid       int       `json:"rid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Testpaper) TableName() string {
	return "testpaper"
}

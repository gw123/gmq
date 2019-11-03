package services

import (
	"time"
)

type ResourceItem struct {
	Id           int32     `json:"id"`
	UserId       int32     `json:"user_id"`
	Content      string    `json:"content"`
	Name         string    `json:"name"`
	Title        string    `json:"title"`
	Avatar       string    `json:"avatar"`
	ChapterTitle string    `json:"chapter_title"`
	GroupTitle   string    `json:"group_title"`
	GroupId      int32     `json:"group_id"`
	ChapterId    int32     `json:"chapter_id"`
	CreatedAt    time.Time `json:"created_at"`
	Type         string    `json:"type"`
}

type GroupItem struct {
	Id       int32       `json:"id"`
	UserId   int32       `json:"user_id"`
	Title    string      `json:"title"`
	Desc     string      `json:"desc"`
	Covers   interface{} `json:"covers"`
	Chapters []*Chapter  `json:"chapters"`
	Tags     []*Tag      `json:"tags"`
	Category Category    `json:"category"`
	UserName string      `json:"user_name"`
	Avatar   string      `json:"avatar"`
}

type Chapter struct {
	Id        int32              `json:"id"`
	Title     string             `json:"title"`
	Resources []*ChapterResource `json:"resources"`
}

type Tag struct {
	Id         int32  `json:"id"`
	Title      string `json:"title"`
	CategoryId int32  `json:"category_id"`
}

func (Tag) TableName() string {
	return "tag"
}

type Category struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
	Tags  []Tag  `json:"tags" gorm:"ForeignKey:category_id"`
}

func (Category) TableName() string {
	return "category"
}

type ChapterResource struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
}

type IndexCtl struct {
	Id      int32       `json:"id"`
	Title   string      `json:"title"`
	Content interface{} `json:"content"`
}

/**
"tpl": "news-item3",
	"data": {
		"title": "\u7ecf\u5178\u7b97\u6cd5\u56fe\u6587\u52a8\u6001\u89e3\u6790",
		"covers": ["http:\/\/data.xytschool.com\/storage\/image\/20190402\/1554214247362272.jpg", "http:\/\/data.xytschool.com\/storage\/image\/20190402\/1554214256699936.jpg", "http:\/\/data.xytschool.com\/storage\/image\/20190402\/1554214269620579.jpg", "http:\/\/data.xytschool.com\/storage\/image\/2019\/04\/05\/a6937643cf28dcbbc3cc6f72c529c33c.jpg", "http:\/\/data.xytschool.com\/storage\/image\/2019\/04\/05\/ac5974c1194c0ed72ee36a3b5fd75a23.png", "http:\/\/data.xytschool.com\/storage\/image\/2019\/04\/05\/1cbbfbf08028ea09473c793a8a3784ba.jpg", "http:\/\/data.xytschool.com\/storage\/image\/2019\/04\/05\/66ec9759b1a5893dde2d163c05f23703."],
		"id": 172,
		"created_at": "2019-03-19 01:51:36",
		"category_id": 40,
		"category": {
			"id": 40,
			"title": "\u7f16\u7a0b",
			"display": null,
			"parent_id": 1,
			"created_at": null,
			"updated_at": "2019-01-23 01:37:40",
			"sort": 0
		},
		"type": "group"
	}
*/
type Block struct {
	Tpl  string `json:"tpl"`
	Type string `json:"type"`
	Data struct {
		Id        int      `json:"id"`
		Title     string   `json:"title"`
		Covers    []string `json:"covers"`
		CreatedAt string   `json:"created_at"`
		Type      string   `json:"type"`
		Url       string   `json:"url"`
		Links     []struct {
			Link  string `json:"link"`
			Title string `json:"title"`
		} `json:"links"`
	} `json:"data"`
}

type ResourceService interface {
	GetServiceName() string
	GetResource(id int) (*ResourceItem, error)
	GetGroupTags(id int32) ([]*Tag, error)
	GetGroupCategory(id int32) (*Category, error)
	GetGroupChapter(id int32) ([]*Chapter, error)
	GetChapterResource(id int32) ([]*ChapterResource, error)
	GetGroup(id int32) (*GroupItem, error)
	GetChapter(id int32) (*Chapter, error)
	GetCategories() ([]*Category, error)
	GetIndexCtrl(maxId, currentId int) ([]*IndexCtl, error)
	GetCategoryCtrl(categoryId, maxId, currentId int) ([]*IndexCtl, error)
	CacheCategoryCtrl(categoryId int) ([]*IndexCtl, error)
	GetCategoryTagCtrl(categoryId, tagId, maxId, currentId int) ([]*GroupItem, error)
	GetQuestions(id int) (*ResourceItem, error)
}

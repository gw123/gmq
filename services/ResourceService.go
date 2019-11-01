package services

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/redisKeys"
	"github.com/gw123/GMQ/common/utils"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
	"strconv"
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

type ResourceService struct {
	app   interfaces.App
	db    *gorm.DB
	redis *redis.Client
}

func NewResourceService(app interfaces.App) (*ResourceService, error) {
	db, err := app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	redisClient, err := app.GetDefaultRedis()
	if err != nil {
		return nil, err
	}
	return &ResourceService{
		app:   app,
		db:    db,
		redis: redisClient,
	}, nil
}

func (s *ResourceService) GetServiceName() string {
	return "ResourceService"
}

func (s *ResourceService) GetResource(id int) (*ResourceItem, error) {
	s.db.LogMode(true)
	var item ResourceItem
	_, err := utils.GetCache(s.redis, redisKeys.Resource+strconv.Itoa(id), &item, func() (interface{}, error) {
		//db.LogMode(true)
		result := s.db.Table("resource as r").
			Select("r.id,r.user_id,r.type,r.title,c.title as chapter_title,c.id as chapter_id ,g.id as group_id,g.title as group_title, r.created_at,users.name,users.avatar").
			Joins("left join users  on r.user_id = users.id").
			//Joins("left join type_article as a on a.rid = r.id").
			Joins("left join chapter as c on c.id = r.chapter_id").
			Joins("left join `group` as g on g.id = r.group_id").
			Where("r.id = ?", id).
			Find(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		var tempItem ResourceItem
		if item.Type == "article" {
			result = s.db.Select("content").Table("type_article").Where("rid = ?", id).Find(&tempItem)
			if result.Error != nil {
				return nil, result.Error
			}
		}else if item.Type == "testpaper" {
			result = s.db.Select("content").Table("testpaper").Where("rid = ?", id).Find(&tempItem)
			if result.Error != nil {
				return nil, result.Error
			}
		}
		item.Content = tempItem.Content
		return item, nil
	})

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ResourceService) GetGroupTags(id int32) ([]*Tag, error) {
	db, err := s.app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	var item []*Tag
	_, err = utils.GetCache(s.redis, redisKeys.GroupTag+strconv.Itoa(int(id)), &item, func() (interface{}, error) {
		var fields = []string{"t.id", "t.title"}
		result := db.Table("`tag` as t").
			Select(fields).
			Joins("left join tag_group as tg  on tg.tag_id = t.id ").
			Joins("left join `group`   as g   on tg.group_id = g.id").
			Where("g.id = ?", id).
			Find(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ResourceService) GetGroupCategory(id int32) (*Category, error) {
	db, err := s.app.GetDefaultDb()
	if err != nil {
		return nil, err
	}
	var item Category
	_, err = utils.GetCache(s.redis, redisKeys.GroupCategory+strconv.Itoa(int(id)), &item, func() (interface{}, error) {
		var fields = []string{"c.id", "c.title"}
		result := db.Table("`group` as g").
			Select(fields).
			Joins("left join `category`  as c on c.id = g.category_id").
			Where("g.id = ?", id).
			Find(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		return &item, nil
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ResourceService) GetGroupChapter(id int32) ([]*Chapter, error) {
	db, err := s.app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	var items []*Chapter
	_, err = utils.GetCache(s.redis, redisKeys.GroupChapter+strconv.Itoa(int(id)), &items, func() (interface{}, error) {
		var fields = []string{"c.id", "c.title"}
		result := db.Table("`chapter` as c").
			Select(fields).
			Where("c.group_id = ?", id).
			Find(&items)
		if result.Error != nil {
			return nil, result.Error
		}
		for i := 0; i < len(items); i++ {
			resources, err := s.GetChapterResource(items[i].Id)
			if err != nil {
				resources = []*ChapterResource{}
			}
			items[i].Resources = resources
		}
		return &items, nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ResourceService) GetChapterResource(id int32) ([]*ChapterResource, error) {
	db, err := s.app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	var item []*ChapterResource
	_, err = utils.GetCache(s.redis, redisKeys.ChapterResource+strconv.Itoa(int(id)), &item, func() (interface{}, error) {
		var fields = []string{"id", "title"}
		result := db.Table("resource").
			Select(fields).
			Where("chapter_id = ?", id).
			Find(&item)
		if result.Error != nil {
			return item, result.Error
		}
		return &item, nil
	})

	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ResourceService) GetGroup(id int32) (*GroupItem, error) {
	db := s.db

	var item GroupItem
	//db.LogMode(true)
	_, err := utils.GetCache(s.redis, redisKeys.Group+strconv.Itoa(int(id)), &item, func() (interface{}, error) {
		var fields = []string{"g.id", "title", "user_id", "`desc`", "covers", "u.name as user_name", "u.avatar"}
		result := db.Table("`group` as g").
			Select(fields).
			Joins("left join users as u on u.id=g.user_id").
			Where("g.id = ?", id).
			Find(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		var covers []string
		err := json.Unmarshal([]byte(item.Covers.([]uint8)), &covers)
		if err != nil {
			s.app.Warn("ResourceService", "json.Unmarshal([]byte(item.Covers.(string)) "+err.Error())
		}
		item.Covers = covers

		chapters, _ := s.GetGroupChapter(id)
		item.Chapters = chapters

		tags, _ := s.GetGroupTags(id)
		item.Tags = tags

		category, err := s.GetGroupCategory(id)
		if err != nil {
			s.app.Warn("ResourceService", "GetGroupCategory err "+err.Error())
			category = &Category{}
		}
		item.Category = *category

		return &item, nil
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *ResourceService) GetChapter(id int32) (*Chapter, error) {
	var item Chapter

	_, err := utils.GetCache(s.redis, redisKeys.Chapter+strconv.Itoa(int(id)), &item, func() (interface{}, error) {
		s.db.Table("chapter").Where("id = ?", id).Find(&item)
		resources, err := s.GetChapterResource(id)
		if err != nil {
			return nil, err
		}
		item.Resources = resources
		return &item, nil
	})

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ResourceService) GetCategories() ([]*Category, error) {
	var items []*Category
	//s.db.LogMode(true)
	_, err := utils.GetCache(s.redis, redisKeys.Categories, &items, func() (interface{}, error) {
		// `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"`
		res := s.db.Where("id != ?", 1).
			Where("is_hide = 0").
			Order("sort").
			Find(&items).Error
		if res != nil {
			return nil, res
		}

		for i := 0; i < len(items); i++ {
			var subs []Tag
			res := s.db.Where("category_id = ?", items[i].Id).Find(&subs).Error
			if res != nil {
				return nil, res
			}
			items[i].Tags = subs
		}
		return &items, nil
	})

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ResourceService) GetIndexCtrl(maxId, currentId int) ([]*IndexCtl, error) {
	var items []*IndexCtl
	//s.db.LogMode(true)
	_, err := utils.GetCache(s.redis, redisKeys.IndexCtrl, &items, func() (interface{}, error) {
		var fields = []string{"i.id", "b.title", "b.content"}

		query := s.db
		if maxId != 0 {
			query.Where("i.id > ?", maxId)
		}
		if currentId != 0 {
			query.Where("i.id < ?", currentId)
		}
		res := query.Table("index_ctrl as i").Select(fields).
			Joins("left join block as b on i.block_id=b.id").
			Order("id desc").
			Find(&items).Error

		if res != nil {
			return nil, res
		}

		for i := 0; i < len(items); i++ {
			var block Block
			err := json.Unmarshal([]byte(items[i].Content.([]uint8)), &block)
			if err != nil {
				return nil, err
			}
			items[i].Content = block
		}

		return &items, nil
	})

	//for i := 0; i < len(items); i++ {
	//	var block Block
	//	buf, err := base64.StdEncoding.DecodeString(items[i].Content.(string))
	//	if err != nil {
	//		return nil, err
	//	}
	//	err = json.Unmarshal(buf, &block)
	//	if err != nil {
	//		return nil, err
	//	}
	//	items[i].Content = block
	//}

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ResourceService) GetCategoryCtrl(categoryId, maxId, currentId int) ([]*IndexCtl, error) {
	var items []*IndexCtl
	var fields = []string{"i.id", "b.title", "b.content"}
	query := s.db.Table("news_ctrl as i").Select(fields).
		Joins("left join block as b on i.block_id=b.id").
		Where("category_id = ?", categoryId).
		Order("id desc").
		Limit(10)
	query.LogMode(true)
	if maxId != 0 {
		query = query.Where("i.id > ?", maxId)
	}
	if currentId != 0 {
		query = query.Where("i.id < ?", currentId)
	}
	res := query.Find(&items).Error
	if res != nil {
		return nil, res
	}

	for i := 0; i < len(items); i++ {
		var block Block
		err := json.Unmarshal([]byte(items[i].Content.([]uint8)), &block)
		if err != nil {
			return nil, err
		}
		items[i].Content = block
	}

	return items, nil
}

func (s *ResourceService) CacheCategoryCtrl(categoryId int) ([]*IndexCtl, error) {
	var items []*IndexCtl
	_, err := utils.GetCache(s.redis, redisKeys.CategoryCtrl+strconv.Itoa(categoryId), &items, func() (interface{}, error) {
		var fields = []string{"i.id", "b.title", "b.content"}
		query := s.db
		res := query.Table("news_ctrl as i").Select(fields).
			Joins("left join block as b on i.block_id=b.id").
			Where("category_id = ?", categoryId).
			Limit(1000).
			Find(&items).Error
		if res != nil {
			return nil, res
		}

		for i := 0; i < len(items); i++ {
			var block Block
			err := json.Unmarshal([]byte(items[i].Content.([]uint8)), &block)
			if err != nil {
				return nil, err
			}
			items[i].Content = block
		}

		return &items, nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ResourceService) GetCategoryTagCtrl(categoryId, tagId, maxId, currentId int) ([]*GroupItem, error) {
	var items []*GroupItem
	_, err := utils.GetCache(s.redis, redisKeys.CategoryCtrl+strconv.Itoa(tagId), &items, func() (interface{}, error) {
		var fields = []string{"g.id", "g.title", "g.desc"}
		query := s.db
		if maxId != 0 {
			query.Where("i.id > ?", maxId)
		}
		if currentId != 0 {
			query.Where("i.id < ?", currentId)
		}
		res := query.Table("`group` as g").Select(fields).
			Joins("left join tag_group as t on t.group_id=g.id").
			Where("t.tag_id = ?", tagId).Limit(50).
			Find(&items).Error
		if res != nil {
			return nil, res
		}

		return &items, nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ResourceService) GetQuestions(id int) (*ResourceItem, error) {
	var item ResourceItem
	_, err := utils.GetCache(s.redis, redisKeys.Resource+strconv.Itoa(id), &item, func() (interface{}, error) {
		//db.LogMode(true)
		result := s.db.Table("resource as r").
			Select("r.id,r.user_id,t.content,r.type,r.title,c.title as chapter_title,c.id as chapter_id ,g.id as group_id,g.title as group_title, r.created_at,users.name,users.avatar").
			Joins("left join users  on r.user_id = users.id").
			Joins("left join testpaper as t on t.rid = r.id").
			Joins("left join chapter as c on c.id = r.chapter_id").
			Joins("left join `group` as g on g.id = r.group_id").
			Where("r.id = ?", id).
			Find(&item)
		if result.Error != nil {
			return nil, result.Error
		}
		return item, nil
	})

	if err != nil {
		return nil, err
	}
	return &item, nil
}

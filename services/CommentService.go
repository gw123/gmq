package services

import (
	"github.com/gw123/GMQ/core/interfaces"
	"time"
)

type CommentItem struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id"`
	Content   string    `json:"content"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentService struct {
	app interfaces.App
}

func NewCommentService(module interfaces.App) *CommentService {
	return &CommentService{
		app: module,
	}
}

func (s *CommentService) GetComments(ctype string, target_id int) ([]*CommentItem, error) {
	db, err := s.app.GetDefaultDb()
	if err != nil {
		return nil, err
	}
	var comments []*CommentItem
	//db.LogMode(true)
	result := db.Table("comments").
		Select("comments.id,user_id,content,comments.created_at,users.name,users.avatar").
		Joins("left join users on comments.user_id = users.id").
		Where("type = ?", ctype).
		Where("target_id = ?", target_id).
		Limit(50).
		Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (s *CommentService) GetServiceName() string {
	return "CommentService"
}

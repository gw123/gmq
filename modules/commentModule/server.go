package commentModule

import (
	"context"
	"github.com/gw123/GMQ/common/models"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type CommentServer struct {
	server   *grpc.Server
	module   interfaces.Module
	bindAddr string
}

func NewCommentServer(module interfaces.Module, bindAddr string) *CommentServer {
	return &CommentServer{
		module:   module,
		bindAddr: bindAddr,
	}
}

func (c *CommentServer) Start() error {
	l, err := net.Listen("tcp", c.bindAddr)
	if err != nil {
		c.module.Error("failed to listen: %v")
		return err
	}
	s := grpc.NewServer() //起一个服务
	RegisterCommentServiceServer(s, c)
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	go s.Serve(l)
	c.module.Info("commentServer listen at : " + c.bindAddr)
	c.server = s
	return nil
}

func (c *CommentServer) Stop() {
	c.server.Stop()
}

func (c *CommentServer) GetComments(ctx context.Context, param *RequestGetComments) (*ResponseGetComments, error) {
	c.module.Info("GetComments request: %v", param)
	if param.TargetId <= 0 {
		c.module.Error("GetComments 参数错误 TargetId 错误")
		return nil, errors.New("参数错误1")
	}

	if !(param.Type == "group" || param.Type == "resource") {
		c.module.Error("GetComments 参数错误 Type 错误")
		return nil, errors.New("参数错误2")
	}

	db, err := c.module.GetApp().GetDefaultDb()
	if err != nil {
		return nil, err
	}

	var comments []*models.Comment
	result := db.Where("type = ?", param.Type).
		Where("target_id = ?", param.TargetId).
		Find(&comments)

	if result.Error != nil {
		return nil, result.Error
	}

	response := ResponseGetComments{
		Comments: make([]*Comment, 0),
	}
	for _, m := range comments {
		c := &Comment{
			Id:        m.ID,
			Type:      m.Type,
			TargetId:  m.TargetId,
			Content:   m.Content,
			UserId:    m.UserId,
			ParentId:  m.ParentId,
			CreatedAt: int32(m.CreatedAt.Unix()),
		}
		response.Comments = append(response.Comments, c)
	}
	return &response, nil
}

func (c *CommentServer) PutComment(ctx context.Context, comment *RequestPutComment) (*ResponsePutComment, error) {
	c.module.Info("PutComment , request: %v", comment)
	commentModel := &models.Comment{
		Type:     comment.Type,
		TargetId: comment.TargetId,
		Content:  comment.GetContent(),
		UserId:   comment.UserId,
		ClientId: comment.ClientId,
		ParentId: comment.ParentId,
	}

	db, err := c.module.GetApp().GetDefaultDb()
	if err != nil {
		return nil, err
	}

	result := db.Save(commentModel)
	if result.Error != nil {
		return nil, result.Error
	}
	response := &ResponsePutComment{Code: 0}
	return response, nil
}

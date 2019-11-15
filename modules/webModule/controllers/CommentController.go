package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gw123/GMQ/common/models"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"strconv"
)

type CommentController struct {
	BaseController
}

func NewCommentController(module interfaces.Module) *CommentController {
	temp := new(CommentController)
	temp.BaseController.module = module
	return temp
}

type commentsParam struct {
}

func (this *CommentController) CommentList(ctx echo.Context) error {
	request := make(map[string]interface{})
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}
	//currentPage := ctx.QueryParam("currentPage")
	//pageSize := ctx.QueryParam("pageSize")

	db, err := this.module.GetApp().GetDefaultDb()
	if err != nil {
		return err
	}
	var comments []*services.CommentItem
	//db.LogMode(true)
	result := db.Table("comments").
		Select("comments.id,user_id,content,comments.created_at,users.name,users.avatar").
		Joins("left join users on comments.user_id = users.id").
		Where("type = ?", request["type"]).
		Where("target_id = ?", request["target_id"]).
		Limit(50).
		Find(&comments)
	if result.Error != nil {
		return result.Error
	}

	return this.BaseController.Success(ctx, comments)
}

func (this *CommentController) Comment(ctx echo.Context) error {
	request := make(map[string]interface{})
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}

	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return this.Fail(ctx, 400, "权限认证失败", errors.New("StandardClaims 类型转换失败"))
	}

	userId, _ := strconv.Atoi(jwtToken.Id)
	commentModel := &models.Comment{
		Ip:       ctx.RealIP(),
		Type:     request["type"].(string),
		TargetId: int32(request["target_id"].(float64)),
		Content:  request["content"].(string),
		UserId:   int32(userId),
		ClientId: request["client_id"].(string),
		//ParentId: request["type"].(string),
	}

	db, err := this.module.GetApp().GetDefaultDb()
	if err != nil {
		return err
	}

	result := db.Save(commentModel)
	if result.Error != nil {
		return result.Error
	}

	return this.BaseController.Success(ctx, nil)
}

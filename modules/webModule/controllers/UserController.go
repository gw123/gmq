package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gw123/GMQ/common/common_types"
	"github.com/gw123/GMQ/common/redisKeys"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type UserController struct {
	BaseController
	UserService *services.UserService
}

func NewUserController(module interfaces.Module) *UserController {
	temp := new(UserController)
	temp.BaseController.module = module
	UserService, ok := module.GetApp().GetService("UserService").(*services.UserService)
	if !ok {
		module.Error("UserService not found")
		return temp
	}
	temp.UserService = UserService
	return temp
}

func (this *UserController) GetUser(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return this.BaseController.Fail(ctx, ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	user, err := this.UserService.GetUser(usrId)
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, user)
}

type SendMessageParam struct {
	Mobile string `json:"mobile"`
}

func (this *UserController) SendMessage(ctx echo.Context) error {
	param := &SendMessageParam{}
	err := ctx.Bind(param)
	if err != nil {
		return this.Fail(ctx, ErrorArgument, "", err)
	}
	randtmp := rand.Int31n(10000)
	if randtmp < 1000 {
		randtmp += 1000
	}
	rand := strconv.Itoa(int(randtmp))

	redis, err := this.module.GetApp().GetDefaultRedis()
	if err != nil {
		return this.Fail(ctx, ErrorDb, "", err)
	}

	if !redis.SetNX(redisKeys.MessageCheckCode+param.Mobile, rand, time.Minute).Val() {
		return this.Fail(ctx, ErrorDb, "请稍等1分钟后尝试", errors.New("请稍等1分钟后尝试"))
	}
	this.module.Info("发送短信:" + param.Mobile)
	event := common_types.NewMobileMessageEvent(rand, param.Mobile)
	this.module.Pub(event)
	return this.BaseController.Success(ctx, nil)
}

func (this *UserController) Register(ctx echo.Context) error {
	params := &services.RegisterParam{}
	err := ctx.Bind(params)
	if err != nil {
		return this.Fail(ctx, ErrorArgument, err.Error(), err)
	}

	err = this.UserService.Register(*params)
	if err != nil {
		return this.Fail(ctx, ErrorDb, err.Error(), err)
	}
	return this.BaseController.Success(ctx, nil)
}

func (this *UserController) Login(ctx echo.Context) error {
	params := &services.LoginParam{}
	err := ctx.Bind(params)
	glog.Dump(params)
	if err != nil {
		return this.Fail(ctx, ErrorArgument, err.Error(), err)
	}

	jwtToken, err := this.UserService.Login(*params)
	if err != nil {
		return this.Fail(ctx, ErrorDb, err.Error(), err)
	}
	data := map[string]string{"api_token": jwtToken}
	glog.Dump(data)
	return this.Success(ctx, data)
}

type ChangeUserCollectionParams struct {
	TargetId     int    `json:"target_id"`
	TargetType   string `json:"target_type"`
	IsCollection int    `json:"isCollect"`
}

func (this *UserController) ChangeUserCollection(ctx echo.Context) error {
	params := &ChangeUserCollectionParams{}
	err := ctx.Bind(params)
	if err != nil {
		return this.Fail(ctx, ErrorArgument, err.Error(), err)
	}

	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return this.Fail(ctx, ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	glog.Dump(params)
	if params.IsCollection == 1 {
		err = this.UserService.AddUserCollection(usrId, params.TargetType, params.TargetId)
		if err != nil {
			return this.Fail(ctx, ErrorDb, err.Error(), err)
		}
	} else {
		err = this.UserService.DelUserCollection(usrId, params.TargetType, params.TargetId)
		if err != nil {
			return this.Fail(ctx, ErrorDb, err.Error(), err)
		}
	}

	return this.Success(ctx, nil)
}

func (this *UserController) GetUserCollection(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return this.Fail(ctx, ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	collection, err := this.UserService.GetUserCollection(usrId)
	if err != nil {
		return this.Fail(ctx, ErrorDb, "", err)
	}

	return this.Success(ctx, collection)
}

package frontend

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/common/redisKeys"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

type UserController struct {
	controllers.BaseController
	UserService *services.UserService
}

func NewUserController(module interfaces.Module) *UserController {
	temp := new(UserController)
	temp.BaseController.Module = module
	UserService, ok := module.GetApp().GetService("UserService").(*services.UserService)
	if !ok {
		module.Error("UserService not found")
		return temp
	}
	temp.UserService = UserService
	return temp
}

func (u *UserController) GetUser(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return u.BaseController.Fail(ctx, controllers.ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	user, err := u.UserService.GetUser(usrId)
	if err != nil {
		return err
	}
	return u.BaseController.Success(ctx, user)
}

type SendMessageParam struct {
	Mobile string `json:"mobile"`
}

func (u *UserController) SendMessage(ctx echo.Context) error {
	param := &SendMessageParam{}
	err := ctx.Bind(param)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorArgument, "", err)
	}
	randtmp := rand.Int31n(10000)
	if randtmp < 1000 {
		randtmp += 1000
	}
	rand := strconv.Itoa(int(randtmp))

	redis, err := u.Module.GetApp().GetDefaultRedis()
	if err != nil {
		return u.Fail(ctx, controllers.ErrorDb, "", err)
	}

	if !redis.SetNX(redisKeys.MessageCheckCode+param.Mobile, rand, time.Minute).Val() {
		return u.Fail(ctx, controllers.ErrorDb, "请稍等1分钟后尝试", errors.New("请稍等1分钟后尝试"))
	}
	u.Module.Info("发送短信:" + param.Mobile)
	event := gmsg.NewMobileMessageEvent(rand, param.Mobile)
	u.Module.Pub(event)
	return u.BaseController.Success(ctx, nil)
}

func (u *UserController) Register(ctx echo.Context) error {
	params := &services.RegisterParam{}
	err := ctx.Bind(params)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorArgument, err.Error(), err)
	}

	err = u.UserService.Register(*params)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorDb, err.Error(), err)
	}
	return u.BaseController.Success(ctx, nil)
}

func (u *UserController) Login(ctx echo.Context) error {
	params := &services.LoginParam{}
	err := ctx.Bind(params)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorArgument, err.Error(), err)
	}

	jwtToken, err := u.UserService.Login(*params)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorDb, err.Error(), err)
	}
	data := map[string]string{"api_token": jwtToken}
	return u.Success(ctx, data)
}

type ChangeUserCollectionParams struct {
	TargetId     int    `json:"target_id"`
	TargetType   string `json:"target_type"`
	IsCollection int    `json:"isCollect"`
}

func (u *UserController) ChangeUserCollection(ctx echo.Context) error {
	params := &ChangeUserCollectionParams{}
	err := ctx.Bind(params)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorArgument, err.Error(), err)
	}

	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return u.Fail(ctx, controllers.ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	if params.IsCollection == 1 {
		err = u.UserService.AddUserCollection(usrId, params.TargetType, params.TargetId)
		if err != nil {
			return u.Fail(ctx, controllers.ErrorDb, err.Error(), err)
		}
	} else {
		err = u.UserService.DelUserCollection(usrId, params.TargetType, params.TargetId)
		if err != nil {
			return u.Fail(ctx, controllers.ErrorDb, err.Error(), err)
		}
	}

	return u.Success(ctx, nil)
}

func (u *UserController) GetUserCollection(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		return u.Fail(ctx, controllers.ErrorArgument, "授权失败", errors.New("jwtToken 解析失败"))
	}
	usrId, _ := strconv.Atoi(jwtToken.Id)
	collection, err := u.UserService.GetUserCollection(usrId)
	if err != nil {
		return u.Fail(ctx, controllers.ErrorDb, "", err)
	}

	return u.Success(ctx, collection)
}

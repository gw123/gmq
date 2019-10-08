package controllers

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"net/http"
)

type IndexController struct {
	BaseController
	UserService *services.UserService
}

func NewIndexController(module interfaces.Module) *IndexController {
	temp := new(IndexController)
	temp.BaseController.module = module
	UserService, ok := module.GetApp().GetService("UserService").(*services.UserService)
	if !ok {
		module.Error("UserService not found")
		return temp
	}
	temp.UserService = UserService
	return temp
}

func (c *IndexController) Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index", nil)
}

func (c *IndexController) Login(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "login", nil)
}

func (c *IndexController) Register(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "register", nil)
}

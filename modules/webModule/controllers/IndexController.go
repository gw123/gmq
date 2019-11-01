package controllers

import (
	"errors"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type IndexController struct {
	BaseController
	UserService     *services.UserService
	ResourceService *services.ResourceService
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

	ResourceService, ok := module.GetApp().GetService("ResourceService").(*services.ResourceService)
	if !ok {
		module.Error("ResourceService not found")
		return temp
	}
	temp.ResourceService = ResourceService
	return temp
}

func (c *IndexController) Index(ctx echo.Context) error {
	currentId, _ := strconv.Atoi(ctx.QueryParam("current_id"))
	maxId, _ := strconv.Atoi(ctx.QueryParam("max_id"))
	items, err := c.ResourceService.GetIndexCtrl(maxId, currentId)
	if err != nil {
		return err
	}
	return ctx.Render(http.StatusOK, "index", map[string]interface{}{"items": items})
}

func (c *IndexController) Group(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	group, err := c.ResourceService.GetGroup(int32(id))
	if err != nil {
		return err
	}
	return ctx.Render(http.StatusOK, "group", map[string]interface{}{"group": group})
}

func (c *IndexController) Resource(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	resource, err := c.ResourceService.GetResource(id)
	if err != nil {
		return err
	}

	if resource.Type == "article" {
		return ctx.Render(http.StatusOK, "resource", map[string]interface{}{"resource": resource})
	} else if resource.Type == "testpaper" {
		return ctx.Render(http.StatusOK, "testpaper", map[string]interface{}{"resource": resource})
	}
	return errors.New("不支持的资源类型")
}

func (c *IndexController) Chapter(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	chapter, err := c.ResourceService.GetChapter(int32(id))
	if err != nil {
		return err
	}
	return ctx.Render(http.StatusOK, "chapter", map[string]interface{}{"chapter": chapter})
}

func (c *IndexController) Home(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "home", nil)
}

func (c *IndexController) News(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "news", nil)
}

func (c *IndexController) TagNews(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "tagNews", nil)
}

func (c *IndexController) Login(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "login", nil)
}

func (c *IndexController) Register(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "register", nil)
}

func (c *IndexController) Testpaper(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	resource, err := c.ResourceService.GetQuestions(id)
	if err != nil {
		return err
	}
	return ctx.Render(http.StatusOK, "testpaper", map[string]interface{}{"resource": resource})
}

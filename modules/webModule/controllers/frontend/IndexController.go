package frontend

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"github.com/skip2/go-qrcode"
	"net/http"
	"strconv"
	"time"
)

type IndexController struct {
	controllers.BaseController
	UserService     *services.UserService
	ResourceService *services.ResourceService
}

func NewIndexController(module interfaces.Module) *IndexController {
	temp := new(IndexController)
	temp.BaseController.Module = module
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
	items, err := c.ResourceService.GetRawIndexCtrl(maxId, currentId)
	if err != nil {
		return err
	}
	return ctx.Render(http.StatusOK, "index", map[string]interface{}{"items": items})
}

func (c *IndexController) Group(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	group, err := c.ResourceService.GetGroup(uint(id))
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
	chapter, err := c.ResourceService.GetChapter(uint(id))
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

func (c *IndexController) Edit(ctx echo.Context) error {
	//id, _ := strconv.Atoi(ctx.Param("id"))
	//resource, err := c.ResourceService.GetQuestions(id)
	//if err != nil {
	//	return err
	//}
	//map[string]interface{}{"resource": resource}
	return ctx.Render(http.StatusOK, "create", nil)
}

//获取一个二维码
func (c *IndexController) Qrcode(ctx echo.Context) error {
	redisClient, err := c.Module.GetApp().GetDefaultRedis()
	if err != nil {
		c.Module.Error("c.Module.GetApp(): " + err.Error())
		return err
	}

	content := ctx.Param("content")
	pngData := make([]byte, 0)
	redisKey := "Qrcode:" + content
	pngData, err = redisClient.Get(redisKey).Bytes()
	if err != nil && err != redis.Nil {
		c.Module.Error("redisClient.get(content): " + err.Error())
	}

	if err == redis.Nil {
		pngData, err = qrcode.Encode(content, qrcode.Highest, 120)
		if err != nil {
			c.Module.Error("qrcode.Encode: " + err.Error())
			return err
		}
		c.Module.Error("not use qrcode cache")
		redisClient.Set(redisKey, pngData, time.Hour)
	}

	count, err := ctx.Response().Write(pngData)
	if err != nil {
		return err
	}
	if count != len(pngData) {
		_, err := ctx.Response().Write(pngData[count:])
		if err != nil {
			return err
		}
	}

	ctx.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(pngData)))
	ctx.Response().Header().Set("Content-Type", "image/png ")
	return nil
}

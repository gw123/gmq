package controllers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gw123/GMQ/common/ctxdata"
	"github.com/gw123/GMQ/common/models"
	"github.com/gw123/GMQ/common/statusCode"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
	"github.com/labstack/echo"
	"strconv"
)

type ResourceController struct {
	BaseController
	ResourceService *services.ResourceService
}

func NewResourceController(module interfaces.Module) *ResourceController {
	temp := new(ResourceController)
	temp.BaseController.module = module
	ResourceService, ok := module.GetApp().GetService("ResourceService").(*services.ResourceService)
	if !ok {
		module.Error("ResourceService not found")
		return temp
	}
	temp.ResourceService = ResourceService
	return temp
}

func (c *ResourceController) GetResource(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		//userId, _ := strconv.Atoi(jwtToken.Id)
	}

	ResourceService, _ := c.module.GetApp().GetService("ResourceService").(*services.ResourceService)
	id, _ := strconv.Atoi(ctx.Param("id"))
	resource, err := ResourceService.GetResource(id)
	if err != nil {
		return err
	}

	return c.BaseController.Success(ctx, resource)
}

func (c *ResourceController) GetGroup(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	group, err := c.ResourceService.GetGroup(uint(id))
	if err != nil {
		return err
	}
	return c.BaseController.Success(ctx, group)
}

func (c *ResourceController) GetChapter(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	chapter, err := c.ResourceService.GetChapter(uint(id))
	if err != nil {
		return err
	}
	return c.BaseController.Success(ctx, chapter)
}

func (c *ResourceController) GetCategories(ctx echo.Context) error {
	items, err := c.ResourceService.GetCategories()
	if err != nil {
		return err
	}
	return c.BaseController.Success(ctx, items)
}

func (c *ResourceController) GetIndexList(ctx echo.Context) error {
	currentId, _ := strconv.Atoi(ctx.QueryParam("current_id"))
	maxId, _ := strconv.Atoi(ctx.QueryParam("max_id"))
	items, err := c.ResourceService.GetIndexCtrl(maxId, currentId)
	if err != nil {
		return err
	}
	return c.BaseController.Success(ctx, items)
}

func (c *ResourceController) GetCategoryCtrl(ctx echo.Context) error {
	categoryId, _ := strconv.Atoi(ctx.Param("category_id"))
	tagId, _ := strconv.Atoi(ctx.Param("tag_id"))
	currentId, _ := strconv.Atoi(ctx.QueryParam("current_id"))
	maxId, _ := strconv.Atoi(ctx.QueryParam("max_id"))
	if tagId != 0 {
		items, err := c.ResourceService.GetCategoryTagCtrl(categoryId, tagId, maxId, currentId)
		if err != nil {
			return err
		}
		return c.BaseController.Success(ctx, items)
	}

	items, err := c.ResourceService.GetCategoryCtrl(categoryId, maxId, currentId)
	if err != nil {
		return err
	}
	return c.BaseController.Success(ctx, items)
}

func (c *ResourceController) SaveGorup(ctx echo.Context) error {
	var group models.Group
	if err := ctx.Bind(&group); err != nil {
		return c.Fail(ctx, ErrorArgument, err.Error(), err)
	}
	if group.ID == 0 {
		group.UserId = ctxdata.GetUserId(ctx)
	} else if !c.ResourceService.CheckGroupAuth(ctx, group) {
		return c.FailCode(ctx, statusCode.NotAuth)
	}

	if err := c.ResourceService.SaveGroup(ctx, &group); err != nil {
		return c.Fail(ctx, ErrorArgument, err.Error(), err)
	}
	return c.Success(ctx, group)
}

func (c *ResourceController) GetRowGroup(ctx echo.Context) error {
	groupIdStr := ctx.Param("id")
	groupId, _ := strconv.Atoi(groupIdStr)
	if groupId == 0 {
		return c.Fail(ctx, ErrorArgument, "请求参数错误", errors.New("请求参数错误"))
	}

	groupItem, err := c.ResourceService.GetRawGroup(uint(groupId))
	if err != nil {
		return c.Fail(ctx, ErrorDb, err.Error(), err)
	}

	if !c.ResourceService.CheckGroupAuth(ctx, groupItem.ID) {
		return c.FailCode(ctx, statusCode.NotAuth)
	}

	return c.Success(ctx, groupItem)
}

func (c *ResourceController) DeleteChapter(ctx echo.Context) error {
	IdStr := ctx.Param("id")
	id, _ := strconv.Atoi(IdStr)
	if id == 0 {
		return c.Fail(ctx, ErrorArgument, "请求参数错误", errors.New("请求参数错误"))
	}

	if !c.ResourceService.CheckChapterAuth(ctx, uint(id)) {
		return c.FailCode(ctx, statusCode.NotAuth)
	}

	err := c.ResourceService.DeleteChapter(uint(id))
	if err != nil {
		return c.Fail(ctx, ErrorDb, err.Error(), err)
	}

	return c.Success(ctx, nil)
}

func (c *ResourceController) SaveResource(ctx echo.Context) error {
	resource := &models.Resource{}
	if err := ctx.Bind(resource); err != nil {
		return c.Fail(ctx, ErrorArgument, err.Error(), err)
	}

	if !c.ResourceService.CheckGroupAuth(ctx, resource.GroupId) {
		return c.FailCode(ctx, statusCode.NotAuth)
	}

	if err := c.ResourceService.SaveResource(ctx, resource); err != nil {
		return c.Fail(ctx, ErrorDb, "保存失败", err)
	}
	return c.Success(ctx, resource)
}

func (c *ResourceController) GetRawResource(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	resource, err := c.ResourceService.GetRawResource(uint(id))
	if err != nil {
		return c.Fail(ctx, ErrorDb, "获取资源失败", err)
	}

	if !c.ResourceService.CheckGroupAuth(ctx, resource.GroupId) {
		return c.FailCode(ctx, statusCode.NotAuth)
	}

	return c.BaseController.Success(ctx, resource)
}

func (c *ResourceController) GetUserGroups(ctx echo.Context) error {
	groups, err := c.ResourceService.GetUserGroups(ctx)
	if err != nil {
		return c.Fail(ctx, ErrorDb, "获取资源失败", err)
	}

	return c.BaseController.Success(ctx, groups)
}

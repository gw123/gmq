package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

func (this *ResourceController) GetResource(ctx echo.Context) error {
	jwtToken, ok := ctx.Get("jwt").(*jwt.StandardClaims)
	if !ok && jwtToken == nil {
		//userId, _ := strconv.Atoi(jwtToken.Id)
	}

	ResourceService, _ := this.module.GetApp().GetService("ResourceService").(*services.ResourceService)
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	resource, err := ResourceService.GetResource(id)
	if err != nil {
		return err
	}

	return this.BaseController.Success(ctx, resource)
}

func (this *ResourceController) GetGroup(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	group, err := this.ResourceService.GetGroup(int32(id))
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, group)
}

func (this *ResourceController) GetChapter(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	chapter, err := this.ResourceService.GetChapter(int32(id))
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, chapter)
}

func (this *ResourceController) GetCategories(ctx echo.Context) error {
	items, err := this.ResourceService.GetCategories()
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, items)
}

func (this *ResourceController) GetIndexList(ctx echo.Context) error {
	currentId, _ := strconv.Atoi(ctx.QueryParam("current_id"))
	maxId, _ := strconv.Atoi(ctx.QueryParam("max_id"))
	items, err := this.ResourceService.GetIndexCtrl(maxId, currentId)
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, items)
}

func (this *ResourceController) GetCategoryCtrl(ctx echo.Context) error {
	categoryId, _ := strconv.Atoi(ctx.Param("category_id"))
	tagId, _ := strconv.Atoi(ctx.Param("tag_id"))
	currentId, _ := strconv.Atoi(ctx.QueryParam("current_id"))
	maxId, _ := strconv.Atoi(ctx.QueryParam("max_id"))
	fmt.Println(currentId)
	if tagId != 0 {
		items, err := this.ResourceService.GetCategoryTagCtrl(categoryId, tagId, maxId, currentId)
		if err != nil {
			return err
		}
		return this.BaseController.Success(ctx, items)
	}

	items, err := this.ResourceService.GetCategoryCtrl(categoryId, maxId, currentId)
	if err != nil {
		return err
	}
	return this.BaseController.Success(ctx, items)
}

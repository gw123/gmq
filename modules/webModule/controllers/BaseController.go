package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"github.com/pkg/errors"
	"github.com/gw123/GMQ/core/interfaces"
	"time"
	"fmt"
	"math/rand"
	"os"
	"io"
)

const (
	Error_NotFound      = 400
	Error_NoAuth        = 401
	Error_DBError       = 402
	Error_CacheError    = 403
	Error_ArgumentError = 404
	Error_NotAllow      = 405
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type BaseController struct {
	module interfaces.Module
}

func NewBaseController(module interfaces.Module) *BaseController {
	baseController := new(BaseController)
	baseController.module = module
	return baseController
}

func (this *BaseController) Success(ctx echo.Context, Data interface{}) error {
	response := webEvent.Response{
		Code: 0 ,
		Msg:  "success",
		Data: Data,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (this *BaseController) Fail(ctx echo.Context, code int, msg string, err error) error {
	response := &webEvent.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	if err == nil {
		err = errors.New(msg)
	}
	err2 := webEvent.NewWebError(response, err)
	return err2
}

func (this *BaseController) uploadFile(ctx echo.Context, formname, cate string) (path string, err error) {
	file, err := ctx.FormFile(formname)
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	path, filename := this.getUploadPath(cate)
	isExist, err := PathExists(path)
	if err != nil {
		return "", err
	}

	if !isExist {
		err := os.MkdirAll(path, 0660)
		if err != nil {
			return "", err
		}
	}
	fullPath := path + "/" + filename
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (this *BaseController) getUploadPath(cate string) (string, string) {
	rootPath := this.module.GetConfig().GetItemOrDefault("uploadRootPath", "./upload")
	pathFormat := this.module.GetConfig().GetItemOrDefault("uploadPathFormat", "mouth")
	relativePath := ""

	switch pathFormat {
	case "day":
		relativePath = time.Now().Format("2019-05-19")
	case "mouth":
		relativePath = time.Now().Format("2019-05")
	}
	num := rand.Int31n(10000)
	filename := fmt.Sprintf("%s_%d", time.Now().String(), num)
	path := rootPath + "/" + "cate/" + relativePath
	return path, filename
}

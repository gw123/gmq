package controllers

import (
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	ErrorNotfound       = 400
	ErrorNotAuth        = 401
	ErrorDb             = 402
	ErrorCache          = 403
	ErrorArgument       = 404
	ErrorNotAllow       = 405
	ErrorEtcd           = 406
	ErrorProtounmarshal = 407
	ErrorProtomarshal   = 408
	ErrorGrpc           = 409
	ErrorTcpWrite       = 410
	ErrorTcpRead        = 411
	ErrorDataTooshort   = 412
	ErrorSignNotmatch   = 413
)

var msgList map[int]string

func init() {
	msgList = make(map[int]string)
	msgList[ErrorNotfound] = "找不对应资源"
	msgList[ErrorNotAuth] = "权限不足"
	msgList[ErrorDb] = "数据库访问出错"
	msgList[ErrorCache] = "缓存访问出错"
	msgList[ErrorArgument] = "请求参数有问题"
	msgList[ErrorNotAllow] = "禁止访问"
	msgList[ErrorEtcd] = "Etcd出错"
	msgList[ErrorProtounmarshal] = "解码数据包出错"
	msgList[ErrorProtomarshal] = "编码数据包出错"
	msgList[ErrorTcpWrite] = "向客户端写入出数据出错"
	msgList[ErrorTcpRead] = "读取客户端数据出错"
	msgList[ErrorDataTooshort] = "发送报文不完整"
	msgList[ErrorSignNotmatch] = "签名校验失败"
}

func ErrorString(code int) string {
	return msgList[code]
}

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
		Code: 0,
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

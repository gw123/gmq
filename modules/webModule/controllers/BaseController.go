package controllers

import (
	"net/http"
	"github.com/labstack/echo"
)

type BaseController struct {
}

func NewBaseController() *BaseController {
	baseController := new(BaseController)
	return baseController
}

func (this *BaseController) Success(ctx echo.Context, Data interface{}) error {
	response := Response{
		Code: 0,
		Msg:  "success",
		Data: Data,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (this *BaseController) Fail(ctx echo.Context, code int, msg string, err error) error {
	response := Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	return ctx.JSON(http.StatusOK, response)
}

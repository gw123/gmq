package controllers

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/labstack/echo"
	"github.com/gw123/GMQ/modules/webModule/models"
)

type ServerController struct {
	BaseController
}

func NewServerController(module interfaces.Module) *ServerController {
	temp := new(ServerController)
	return temp
}

/***
 * 登录并且检测是否有新的版本
 */
func (this *ServerController) Login(ctx echo.Context) error {
	version := ctx.QueryParam("version")
	formData := make(map[string]string, 0)
	formData["version"] = version

	return this.Success(ctx, formData)
}

/***
 * 上传新的版本
 */
func (this *ServerController) UploadVersion(ctx echo.Context) error {
	server := new(models.Service)
	if err := ctx.Bind(server); err != nil {
		return this.Fail(ctx, 0, "参数解析失败" ,err)
	}

	if err := ctx.Validate(server); err != nil {
		return err
	}
	//fmt.Println(server)
	//formData := make(map[string]string,0)
	//formData["version"] = version
	return this.Success(ctx, server)
}

/***
 * 下载新的版本
 */
func (this *ServerController) Download(ctx echo.Context) error {

	return this.Success(ctx, nil)
}

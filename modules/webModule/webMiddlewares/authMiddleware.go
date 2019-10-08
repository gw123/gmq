package webMiddlewares

import (
	"github.com/gw123/GMQ/common/utils"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func NewAuthMiddleware(module interfaces.Module) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if ctx.Request().Method == http.MethodOptions {
				return ctx.HTML(http.StatusOK, "")
			}
			authToen := ctx.Request().Header.Get("Authorization")
			temps := strings.Split(authToen, " ")
			if len(temps) == 2 {
				authToen = temps[1]
			} else {
				ret := make(map[string]interface{}, 0)
				ret["code"] = 401
				ret["msg"] = "授权失败"
				return ctx.JSON(http.StatusOK, ret)
			}

			jwt_key := module.GetConfig().GetStringItem("jwt_key")

			infos, err := utils.ParseJwsTokenSh1(authToen, jwt_key)
			if err != nil {
				ret := make(map[string]interface{}, 0)
				ret["code"] = 401
				ret["msg"] = "授权失败1:" + err.Error()
				return ctx.JSON(http.StatusOK, ret)
			}
			ctx.Set("jwt", infos)
			return next(ctx)
		}
	}
}

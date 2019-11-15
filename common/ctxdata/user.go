package ctxdata

import (
	"github.com/labstack/echo"
)

const _userId = "user_id"

func SetUserId(ctx echo.Context, userId uint) {
	ctx.Set(_userId, userId)
}

func GetUserId(ctx echo.Context) uint {
	temp := ctx.Get(_userId)
	id, _ := temp.(uint)
	return id
}

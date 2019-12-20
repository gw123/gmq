package gmqcore

import (
	"github.com/gw123/gmq"
)

type MiddlewareManager struct {
	middlewares []gmq.Middleware
	app         gmq.App
}

func NewMiddlewareManager(app gmq.App) *MiddlewareManager {
	this := new(MiddlewareManager)
	this.app = app
	return this
}

func (this *MiddlewareManager) Handel(event gmq.Msg) bool {
	for _, middleware := range this.middlewares {
		flag := middleware.Handel(event)
		if !flag {
			return false
		}
	}
	return true
}

func (this *MiddlewareManager) RegisterMiddleware(middleware gmq.Middleware) {
	this.middlewares = append(this.middlewares, middleware)
}

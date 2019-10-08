package bootstarp

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/services"
)

func LoadServices(app interfaces.App) {
	comment := services.NewCommentService(app)
	app.RegisterService("", comment)

	resource, err := services.NewResourceService(app)
	if err != nil {
		panic(err)
	}
	app.RegisterService("", resource)

	userService, err := services.NewUserService(app)
	if err != nil {
		panic(err)
	}
	app.RegisterService("", userService)
	return
}

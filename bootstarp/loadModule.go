package bootstarp

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/MessageModule"
	redisModule "github.com/gw123/GMQ/modules/RedisModule"
	"github.com/gw123/GMQ/modules/commentModule"
	"github.com/gw123/GMQ/modules/debugModule"
	"github.com/gw123/GMQ/modules/webModule"
)

func LoadModuleProvider(app interfaces.App) {
	app.LoadModuleProvider(debugModule.NewDebugModuleProvider())
	app.LoadModuleProvider(webModule.NewWebModuleProvider())
	app.LoadModuleProvider(MessageModule.NewMessageModuleProvider())
	app.LoadModuleProvider(redisModule.NewRedisModuleProvider())
	//app.LoadModuleProvider(serverNodeModule.NewServerNodeModuleProvider())
	app.LoadModuleProvider(commentModule.NewDebugModuleProvider())
	return
}

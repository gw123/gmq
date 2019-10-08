package test

import (
	"github.com/gw123/GMQ/bootstarp"
	"github.com/gw123/GMQ/core"
	"github.com/gw123/GMQ/core/interfaces"
)

func GetAppFroTest() interfaces.App {
	bootstarp.SetConfigFile("./config.yml")
	config := bootstarp.GetConfig()
	App := core.NewApp(config)
	App.Start()
	return App
}

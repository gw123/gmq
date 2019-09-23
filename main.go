package main

import (
	"flag"
	"github.com/gw123/GMQ/bootstarp"
	"github.com/gw123/GMQ/core"
)

func parseConfig() {
	configFile := flag.String("c", "config.yml", "配置文件")
	flag.Parse()
	bootstarp.SetConfigFile(*configFile)
}

func main() {
	parseConfig()
	config := bootstarp.GetConfig()
	appInstance := core.NewApp(config)
	bootstarp.LoadModuleProvider(appInstance)
	//migrage tabels   迁移数据库
	err := bootstarp.AutoMigrate(appInstance)
	if err != nil {
		appInstance.Error("App", err.Error())
	}
	appInstance.Start()
	select {}
}

package main

import (
	"github.com/gw123/GMQ/core"
	"github.com/gw123/GMQ/bootstarp"
	"flag"
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
	appInstance.Start()

	select {}
}

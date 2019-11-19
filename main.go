package main

import (
	"flag"
	"fmt"
	"github.com/gw123/GMQ/bootstarp"
	"github.com/gw123/GMQ/common/helper"
	"github.com/gw123/GMQ/core"
	//"net/http"
	//_ "net/http/pprof"
)

func main() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8282", nil))
	//}()

	//defer profile.Start().Stop()

	moduleName := flag.String("newModule", "", "模块名称  eg:RedisModule")
	serviceName := flag.String("newService", "", "service名称  eg:UserService")
	configFile := flag.String("c", "config.yml", "配置文件")
	dir := flag.String("dir", "", "输出位置  默认在当前的modules下面")

	flag.Parse()

	if *moduleName != "" {
		err := helper.MakeModule(*moduleName, *dir)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("created " + *moduleName + " success!")
		}
		return
	}

	if *serviceName != "" {
		err := helper.MakeService(*serviceName, *dir)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("created " + *serviceName + " success!")
		}
		return
	}
	bootstarp.SetConfigFile(*configFile)

	config := bootstarp.GetConfig()
	appInstance := core.NewApp(config)

	//migrage tabels   迁移数据库
	err := bootstarp.AutoMigrate(appInstance)
	if err != nil {
		appInstance.Error("App", err.Error())
	}

	//LoadServices 加载服务
	bootstarp.LoadServices(appInstance)
	//LoadModuleProvider 加载模块
	bootstarp.LoadModuleProvider(appInstance)

	appInstance.Start()
	select {}
}

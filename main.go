package main

import (
	_ "net/http/pprof"
	"github.com/gw123/GMQ/core"
	"github.com/gw123/GMQ/bootstarp"
)

func main() {
	config := bootstarp.GetConfig()
	appInstance := core.NewApp(config)
	bootstarp.LoadModuleProvider(appInstance)
	appInstance.Start()
	//go func() {
	//	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	//}()

	select {
	}
}

// + !debug

package app

import (
	"github.com/gw123/GMQ/modules/debugModule"
	"github.com/gw123/GMQ/modules/mqttModule"
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/interfaces"
	"errors"
	"fmt"
	"github.com/fpay/erp-client-s/modules/printerModule"
)

/*
 * 模块管理模块
 * 1. 加载模块,卸载模块
 * 2. 管理配置,更新模块配置
*/
type ModuleManager struct {
	configManager *ConfigManager
	Modules       map[string]interfaces.Module
	app           interfaces.App
	ConfigData    []byte
}

func NewModuleManager(app interfaces.App, configManger *ConfigManager) *ModuleManager {
	this := new(ModuleManager)
	this.app = app
	this.configManager = configManger
	this.Modules = make(map[string]interfaces.Module)
	return this
}

func (this *ModuleManager) LoadModules() {
	for moduleName, moduleConfig := range this.configManager.ModuleConfigs {
		if moduleConfig.IsEnable() == false {
			this.app.Info("ConfigManager", "禁止加载 "+moduleName)
			continue
		}
		err := this.LoadModule(moduleName, moduleConfig)
		if err != nil {
			this.app.Error("ConfigManager", "模块加载失败 "+moduleName+" "+err.Error())
		} else {
			this.app.Info("ConfigManager", "加载成功 "+moduleName)
		}
	}
}

func (this *ModuleManager) LoadModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	if config.IsInnerModule() {
		this.app.Info("ConfigManager", "加载内部模块 "+moduleName)
		err = this.loadInnerModule(moduleName, config)
	} else {
		this.app.Info("ConfigManager", "加载外部模块 "+moduleName)
		err = this.loadDll(moduleName, config)
	}
	return
}

func (this *ModuleManager) UnLoadModule(moduleName string) (err error) {
	if this.Modules[moduleName] == nil {
		return
	}
	this.Modules[moduleName].UnInit()
	return
}

func (this *ModuleManager) loadDll(muduleName string, config interfaces.ModuleConfig) (err error) {
	module := base.NewDllModule()
	err = module.Init(this.app, config)
	return err
}

func (this *ModuleManager) loadInnerModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	switch moduleName {
	case "PrinterModule":
		this.Modules[moduleName] = printerModule.NewPrinterModule()
		err = this.Modules[moduleName].Init(this.app, config)
		break
	case "MqttModule":
		this.Modules[moduleName] = mqttModule.NewMqttModule()
		err = this.Modules[moduleName].Init(this.app, config)
		break
	case "DebugModule":
		this.Modules[moduleName] = debugModule.NewDebugModule()
		err = this.Modules[moduleName].Init(this.app, config)
		break
		//case "WebSocketModule":
		//	this.Modules[moduleName] = webSocketModule.NewWebSocketModule()
		//	err = this.Modules[moduleName].Init(this.app, config)
		//	break;
	default:
		err = errors.New("没有这样的模块")
	}
	return err
}

func (this ModuleManager) GetModuleStatus() string {
	str := ""
	for moduleName, module := range this.Modules {
		str += fmt.Sprintf("%: %\n", moduleName, module.GetEventNum())
	}
	return str
}

//func (this *ConfigManager) LoadModule1(muduleName string, config []byte) {
//	module := this.Modules[muduleName]
//	module.Init(this.AppInstance, config)
//}

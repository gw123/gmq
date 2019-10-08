// + !debug

package core

import (
	"errors"
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
	"strings"
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
	providers     map[string]interfaces.ModuleProvider
}

func NewModuleManager(app interfaces.App, configManger *ConfigManager) *ModuleManager {
	this := new(ModuleManager)
	this.app = app
	this.configManager = configManger
	this.Modules = make(map[string]interfaces.Module)
	this.providers = make(map[string]interfaces.ModuleProvider, 0)
	return this
}

func (this *ModuleManager) LoadModules() {
	for moduleName, moduleConfig := range this.configManager.ModuleConfigs {
		if moduleConfig.IsEnable() == false {
			//this.app.Info("ConfigManager", "禁止加载 "+moduleName)
			continue
		}
		err := this.LoadModule(moduleName, moduleConfig)
		if err != nil {
			this.app.Error("ModuleManager", "模块加载失败 on init"+moduleName+" "+err.Error())
		} else {
			this.app.Info("ModuleManager", "加载成功 "+moduleName)
		}

		err = this.Modules[moduleName].BeforeStart()
		if err != nil {
			this.app.Error("ModuleManager", "模块加载失败 on BeforeStart "+moduleName+" "+err.Error())
		} else {
			this.app.Info("ModuleManager", "加载成功 "+moduleName)
		}

	}

	for _, module := range this.Modules {
		if module != nil {
			go module.Start()
		}
	}
}

func (this *ModuleManager) LoadModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	moduleType := config.GetModuleType()
	if moduleType == "" {
		moduleType = "inner"
	}

	switch moduleType {
	case "inner":
		this.app.Info("ModuleManager", "加载内部模块 "+moduleName)
		err = this.loadInnerModule(moduleName, config)
		break
	case "dll":
	case "lib":
		this.app.Info("ModuleManager", "加载外部模块 "+moduleName)
		err = this.loadDll(moduleName, config)
		break
	case "exe":
		this.app.Info("ModuleManager", "加载外部模块 "+moduleName)
		err = this.loadExe(moduleName, config)
		break
	default:
		err = errors.New("not support module type 不支持的模块类型")
	}
	return err
}

func (this *ModuleManager) UnLoadModule(moduleName string) (err error) {
	if this.Modules[moduleName] == nil {
		return
	}
	this.Modules[moduleName].UnInit()
	return
}

func (this *ModuleManager) loadDll(moduleName string, config interfaces.ModuleConfig) (err error) {
	module := base.NewDllModule()
	err = module.Init(this.app, config)

	if err == nil {
		this.Modules[moduleName] = module
	}

	return err
}

func (this *ModuleManager) loadExe(moduleName string, config interfaces.ModuleConfig) (err error) {
	module := base.NewExeModule()
	err = module.Init(this.app, config)
	if err == nil {
		this.Modules[moduleName] = module
	}
	return
}

//注意模块统一小写
func (this *ModuleManager) LoadModuleProvider(provider interfaces.ModuleProvider) {
	if provider == nil {
		return
	}
	this.providers[strings.ToLower(provider.GetModuleName())] = provider
}

func (this *ModuleManager) loadInnerModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	provider, ok := this.providers[moduleName]
	if ok {
		newModule := provider.GetModule()
		err = newModule.Init(this.app, config)
		if err == nil {
			this.Modules[moduleName] = newModule
		}
	} else {
		err = errors.New("没有这样的模块")
	}
	return err
}

func (this ModuleManager) GetModuleStatus() string {
	str := ""
	for moduleName, module := range this.Modules {
		str += fmt.Sprintf("%s: %d\n", moduleName, module.GetStatus())
	}
	return str
}

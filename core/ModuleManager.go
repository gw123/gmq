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

func (m *ModuleManager) LoadModules() {
	for moduleName, moduleConfig := range m.configManager.ModuleConfigs {
		if moduleConfig.IsEnable() == false {
			//m.app.Debug("ConfigManager","enable:%v,", moduleConfig.GetItem("enable"))
			m.app.Info("ConfigManager", "禁止加载 "+moduleName)
			continue
		}
		err := m.LoadModule(moduleName, moduleConfig)
		if err != nil {
			m.app.Error("ModuleManager", "模块加载失败 "+moduleName+" "+err.Error())
			continue
		} else {
			m.app.Info("ModuleManager", "加载成功 "+moduleName)
		}

		err = m.Modules[moduleName].BeforeStart()
		if err != nil {
			m.app.Error("ModuleManager", "模块加载失败 on BeforeStart "+moduleName+" "+err.Error())
			continue
		} else {
			m.app.Info("ModuleManager", "加载成功 "+moduleName)
		}
	}

	for _, module := range m.Modules {
		if module != nil {
			go module.Start()
		}
	}
}

func (m *ModuleManager) LoadModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	moduleType := config.GetModuleType()
	if moduleType == "" {
		moduleType = "inner"
	}

	switch moduleType {
	case "inner":
		m.app.Info("ModuleManager", "加载内部模块 "+moduleName)
		err = m.loadInnerModule(moduleName, config)
		break
	case "dll":
	case "lib":
		m.app.Info("ModuleManager", "加载外部模块 "+moduleName)
		err = m.loadDll(moduleName, config)
		break
	case "exe":
		m.app.Info("ModuleManager", "加载外部模块 "+moduleName)
		err = m.loadExe(moduleName, config)
		break
	default:
		err = errors.New("not support module type 不支持的模块类型")
	}
	return err
}

func (m *ModuleManager) UnLoadModule(moduleName string) (err error) {
	if m.Modules[moduleName] == nil {
		return
	}
	m.Modules[moduleName].UnInit()
	return
}

func (m *ModuleManager) loadDll(moduleName string, config interfaces.ModuleConfig) (err error) {
	module := base.NewDllModule()
	err = module.Init(m.app, config)

	if err == nil {
		m.Modules[moduleName] = module
	}

	return err
}

func (m *ModuleManager) loadExe(moduleName string, config interfaces.ModuleConfig) (err error) {
	module := base.NewExeModule()
	err = module.Init(m.app, config)
	if err == nil {
		m.Modules[moduleName] = module
	}
	return
}

//注意模块统一小写
func (m *ModuleManager) LoadModuleProvider(provider interfaces.ModuleProvider) {
	if provider == nil {
		return
	}
	m.providers[strings.ToLower(provider.GetModuleName())] = provider
}

func (m *ModuleManager) loadInnerModule(moduleName string, config interfaces.ModuleConfig) (err error) {
	provider, ok := m.providers[moduleName]
	if ok {
		newModule := provider.GetModule()
		err = newModule.Init(m.app, config)
		if err == nil {
			m.Modules[moduleName] = newModule
		}
	} else {
		err = errors.New("没有这样的模块")
	}
	return err
}

func (m ModuleManager) GetModuleStatus() string {
	str := ""
	for moduleName, module := range m.Modules {
		str += fmt.Sprintf("%s: %d\n", moduleName, module.GetStatus())
	}
	return str
}

package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"sync"
	"fmt"
)

type ModuleConfig struct {
	mutex        sync.Mutex
	ModuleName   string                 `json:"module_name"`
	Configs      map[string]interface{} `json:"configs"`
	GlobalConfig interfaces.AppConfig
}

func NewModuleConfig(moduleName string, appConfig interfaces.AppConfig) *ModuleConfig {
	this := new(ModuleConfig)
	this.ModuleName = moduleName
	this.Configs = make(map[string]interface{})
	this.GlobalConfig = appConfig
	this.SetItem("subs", "")
	return this
}

func (this *ModuleConfig) GetPath() string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	value, ok := this.Configs["path"].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败", this.GetModuleName(), "path")
	return value
}

//默认情况下是内部模块
func (this *ModuleConfig) IsInnerModule() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.Configs["inner"] == "true" ||
		this.Configs["inner"] == "1" ||
		this.Configs["inner"] == "" ||
		this.Configs["inner"] == nil
}

//是否启动
func (this *ModuleConfig) IsEnable() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.Configs["enable"] == "true" ||
		this.Configs["enable"] == "1" ||
		this.Configs["enable"] == "" ||
		this.Configs["enable"] == nil
}

func (this *ModuleConfig) GetModuleName() string {
	return this.ModuleName
}

func (this *ModuleConfig) SetItem(key string, value interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.Configs[key] = value
}

func (this *ModuleConfig) GetItems() (value map[string]interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.Configs
}

func (this *ModuleConfig) GetGlobalItem(key string) (value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.GlobalConfig.GetItem(key)
}

func (this *ModuleConfig) GetGlobalItems() (value map[string]interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.GlobalConfig.GetItems()
}

func (this *ModuleConfig) SetGlobalConfig(config interfaces.AppConfig) {
	if config == nil {
		return
	}
	this.GlobalConfig = config
}

func (this *ModuleConfig) GetItem(key string) (value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	value, ok := this.Configs[key].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return ""
}

func (this *ModuleConfig) GetIntItem(key string) int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	value, ok := this.Configs[key].(int)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return 0
}

func (this *ModuleConfig) GetBoolItem(key string) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	value, ok := this.Configs[key].(bool)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return false
}

func (this *ModuleConfig) MergeNewConfig(newCofig interfaces.ModuleConfig) bool {
	newConfigItems := newCofig.GetGlobalItems()
	isChange := false
	for key, newConfigItem := range newConfigItems {
		oldConfig := this.GetItem(key)
		if newConfigItem != oldConfig {
			this.SetItem(key, newConfigItem)
			isChange = true
		}
	}
	return isChange
}
func (this *ModuleConfig) GetModuleType() string {
	return this.GetItem("type")
}

package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"sync"
)

type ModuleConfig struct {
	sync.Mutex
	moduleName   string
	configs      map[string]string
	GlobalConfig interfaces.AppConfig
}

func NewModuleConfig(moduleName string, appConfig interfaces.AppConfig) *ModuleConfig {
	this := new(ModuleConfig)
	this.moduleName = moduleName
	this.configs = make(map[string]string)
	this.GlobalConfig = appConfig
	return this
}

func (this *ModuleConfig) GetPath() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs["path"]
}

//默认情况下是内部模块
func (this *ModuleConfig) GetModuleType() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs["type"]
}

//是否启动
func (this *ModuleConfig) IsEnable() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs["enable"] == "true" || this.configs["enable"] == "1" || this.configs["enable"] == ""
}

func (this *ModuleConfig) GetModuleName() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.moduleName
}

func (this *ModuleConfig) GetItem(key string) (value string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs[key]
}

func (this *ModuleConfig) SetItem(key, value string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.configs[key] = value
}

func (this *ModuleConfig) GetItems() (value map[string]string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs
}

func (this *ModuleConfig) GetGlobalItem(key string) (value string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.GlobalConfig.GetItem(key)
}

func (this *ModuleConfig) GetGlobalItems() (value map[string]string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.GlobalConfig.GetItems()
}

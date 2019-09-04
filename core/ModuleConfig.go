package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"sync"
	"fmt"
	"strings"
)

type ModuleConfig struct {
	mutex        sync.RWMutex
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
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	value, ok := this.Configs["path"].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败", this.GetModuleName(), "path")
	return value
}

//默认情况下是内部模块
func (this *ModuleConfig) IsInnerModule() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return this.Configs["inner"] == "true" ||
		this.Configs["inner"] == "1" ||
		this.Configs["inner"] == "" ||
		this.Configs["inner"] == nil
}

//是否启动
func (this *ModuleConfig) IsEnable() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.Configs["enable"] == "true" ||
		this.Configs["enable"] == "1" ||
		this.Configs["enable"] == "" ||
		this.Configs["enable"] == nil
}

func (this *ModuleConfig) GetModuleName() string {
	return this.ModuleName
}

func (this *ModuleConfig) SetItem(key string, value interface{}) {
	key1 := strings.ToLower(key)
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.Configs[key1] = value
}

func (this *ModuleConfig) GetItems() (value map[string]interface{}) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.Configs
}

func (this *ModuleConfig) GetGlobalItem(key string) (value string) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.GlobalConfig.GetItem(key)
}

func (this *ModuleConfig) GetGlobalItems() (value map[string]interface{}) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.GlobalConfig.GetItems()
}

func (this *ModuleConfig) SetGlobalConfig(config interfaces.AppConfig) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if config == nil {
		return
	}
	this.GlobalConfig = config
}

func (this *ModuleConfig) GetItem(key string) (value string) {
	key1 := strings.ToLower(key)
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key1].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return ""
}

func (this *ModuleConfig) GetArrayItem(key string) (value []string) {
	key1 := strings.ToLower(key)
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	//fmt.Println(this.Configs[key1])
	values, ok := this.Configs[key1].([]interface{})
	if !ok {
		fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
		return nil
	}

	value = make([]string, 0)
	for _, val := range values {
		if r, ok := val.(string); ok {
			value = append(value, r)
		}
	}
	return value
}

func (this *ModuleConfig) GetMapItem(key string) (value map[string]string) {
	key1 := strings.ToLower(key)
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key1].(map[string]string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return nil
}

func (this *ModuleConfig) GetIntItem(key string) int {
	key1 := strings.ToLower(key)
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key1].(int)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", this.GetModuleName(), key)
	return 0
}

func (this *ModuleConfig) GetBoolItem(key string) bool {
	key1 := strings.ToLower(key)
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key1].(bool)
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

func (this *ModuleConfig) GetItemOrDefault(key, defaultval string) string {
	ret := this.GetItem(key)
	if ret == "" {
		return defaultval
	}
	return ret
}

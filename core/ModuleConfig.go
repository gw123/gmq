package core

import (
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"strings"
	"sync"
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

func (m *ModuleConfig) GetPath() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, ok := m.Configs["path"].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败", m.GetModuleName(), "path")
	return value
}

//默认情况下是内部模块
func (m *ModuleConfig) IsInnerModule() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.Configs["inner"] == "true" ||
		m.Configs["inner"] == "1" ||
		m.Configs["inner"] == "" ||
		m.Configs["inner"] == nil
}

//是否启动
func (m *ModuleConfig) IsEnable() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Configs["enable"] == "true" ||
		m.Configs["enable"] == true ||
		m.Configs["enable"] == "1" ||
		m.Configs["enable"] == "" ||
		m.Configs["enable"] == "yes" ||
		m.Configs["enable"] == nil
}

func (m *ModuleConfig) GetModuleName() string {
	return m.ModuleName
}

func (m *ModuleConfig) SetItem(key string, value interface{}) {
	key1 := strings.ToLower(key)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Configs[key1] = value
}

func (m *ModuleConfig) GetItems() (value map[string]interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Configs
}

func (m *ModuleConfig) GetItem(key string) (value interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.Configs[key]
}

func (m *ModuleConfig) GetGlobalItem(key string) (value string) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.GlobalConfig.GetItem(key)
}

func (m *ModuleConfig) GetGlobalItems() (value map[string]interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.GlobalConfig.GetItems()
}

func (m *ModuleConfig) SetGlobalConfig(config interfaces.AppConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if config == nil {
		return
	}
	m.GlobalConfig = config
}

func (m *ModuleConfig) GetStringItem(key string) (value string) {
	key1 := strings.ToLower(key)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Configs[key1].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", m.GetModuleName(), key)
	return ""
}

func (m *ModuleConfig) GetArrayItem(key string) (value []string) {
	key1 := strings.ToLower(key)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	//fmt.Println(m.Configs[key1])
	values, ok := m.Configs[key1].([]interface{})
	if !ok {
		fmt.Printf("模块 %s 获取配置 %s 失败\n", m.GetModuleName(), key)
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

func (m *ModuleConfig) GetMapItem(key string) (value map[string]interface{}) {
	key1 := strings.ToLower(key)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Configs[key1].(map[string]interface{})
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", m.GetModuleName(), key)
	return nil
}

func (m *ModuleConfig) GetIntItem(key string) int {
	key1 := strings.ToLower(key)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Configs[key1].(int)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", m.GetModuleName(), key)
	return 0
}

func (m *ModuleConfig) GetBoolItem(key string) bool {
	key1 := strings.ToLower(key)
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Configs[key1].(bool)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", m.GetModuleName(), key)
	return false
}

func (m *ModuleConfig) MergeNewConfig(newCofig interfaces.ModuleConfig) bool {
	newConfigItems := newCofig.GetGlobalItems()
	isChange := false
	for key, newConfigItem := range newConfigItems {
		oldConfig := m.GetStringItem(key)
		if newConfigItem != oldConfig {
			m.SetItem(key, newConfigItem)
			isChange = true
		}
	}
	return isChange
}
func (m *ModuleConfig) GetModuleType() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.GetStringItem("type")
}

func (m *ModuleConfig) GetItemOrDefault(key, defaultval string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	ret := m.GetStringItem(key)
	if ret == "" {
		return defaultval
	}
	return ret
}

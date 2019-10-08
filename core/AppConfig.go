package core

import (
	"sync"
	"fmt"
)

type AppConfig struct {
	mutex   sync.RWMutex
	Configs map[string]interface{}
}

func NewAppConfig() *AppConfig {
	this := new(AppConfig)
	this.Configs = make(map[string]interface{})
	return this
}

func (this *AppConfig) GetItem(key string) (value string) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key].(string)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", "App", key)
	return ""
}

func (this *AppConfig) GetIntItem(key string) int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key].(int)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", "App", key)
	return 0
}

func (this *AppConfig) GetBoolItem(key string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	value, ok := this.Configs[key].(bool)
	if ok {
		return value
	}
	fmt.Printf("模块 %s 获取配置 %s 失败\n", "App", key)
	return false
}

func (this *AppConfig) SetItem(key, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.Configs[key] = value
}

func (this *AppConfig) GetItems() (value map[string]interface{}) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.Configs
}

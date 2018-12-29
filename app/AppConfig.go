package app

import "sync"

type AppConfig struct {
	sync.Mutex
	configs    map[string]string
}

func NewAppConfig() *AppConfig {
	this := new(AppConfig)
	this.configs = make(map[string]string)
	return this
}

func (this *AppConfig) GetItem(key string) (value string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs[key]
}

func (this *AppConfig) SetItem(key, value string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.configs[key] = value
}

func (this *AppConfig) GetItems() (value map[string]string) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.configs
}

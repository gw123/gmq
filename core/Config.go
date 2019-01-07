// + !debug

package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/go-ini/ini"
	"fmt"
)

/*
 * 模块管理模块
 * 1. 加载模块,卸载模块
 * 2. 管理配置,更新模块配置
*/
type ConfigManager struct {
	ModuleConfigs map[string]interfaces.ModuleConfig
	app           interfaces.App
	GlobalConfig  interfaces.AppConfig
	ConfigData    []byte
}

func NewConfigManager(app interfaces.App, configData []byte) *ConfigManager {
	this := new(ConfigManager)
	this.app = app
	this.ModuleConfigs = make(map[string]interfaces.ModuleConfig)
	this.GlobalConfig = NewAppConfig()
	this.ConfigData = configData
	err := this.ParseConfig()
	if err != nil {
		this.app.Warning("ConfigManger", "配置文件解析失败 "+err.Error())
	}
	return this
}

func (this *ConfigManager) ParseConfig() (err error) {
	cfg, err := ini.Load(this.ConfigData)
	if err != nil {
		return
	}
	sections := cfg.Sections()
	section := cfg.Section("DEFAULT")
	keys := section.Keys()
	for _, key := range keys {
		//fmt.Println("\t", key.Name(), ":", key.String())
		this.GlobalConfig.SetItem(key.Name(), key.String())
	}

	for _, section := range sections {
		if section.Name() == "DEFAULT" {
			continue
		}
		//fmt.Println("section name :", section.Name())
		moduleConfig := NewModuleConfig(section.Name(), this.GlobalConfig)
		keys := section.Keys()
		for _, key := range keys {
			//fmt.Println("\t", key.Name(), ":", key.String())
			moduleConfig.SetItem(key.Name(), key.String())
		}
		this.ModuleConfigs[section.Name()] = moduleConfig
	}
	return
}

func (this *ConfigManager) ParseJsonConfig(configData string) (err error) {
	cfg, err := ini.Load(configData)
	if err != nil {
		fmt.Printf("Fail to load %v", err)
		return
	}
	sections := cfg.Sections()
	section := cfg.Section("DEFAULT")
	keys := section.Keys()
	for _, key := range keys {
		//fmt.Println("\t", key.Name(), ":", key.String())
		this.GlobalConfig.SetItem(key.Name(), key.String())
	}

	for _, section := range sections {
		if section.Name() == "DEFAULT" {
			continue
		}
		//fmt.Println("section name :", section.Name())
		moduleConfig := NewModuleConfig(section.Name(), this.GlobalConfig)
		keys := section.Keys()
		for _, key := range keys {
			//fmt.Println("\t", key.Name(), ":", key.String())
			moduleConfig.SetItem(key.Name(), key.String())
		}
		this.ModuleConfigs[section.Name()] = moduleConfig
	}
	return
}

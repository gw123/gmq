package core

import (
	"github.com/gw123/GMQ/common/common_types"
	"github.com/gw123/GMQ/core/middlewares"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/go-ini/ini"
	"encoding/json"
	"fmt"
)

type App struct {
	configFile        *ini.File
	errorManager      *ErrorManager
	moduleManager     *ModuleManager
	configManager     *ConfigManager
	middlewareManager *MiddlewareManager
	logManager        *LogManager
	dispatch          *Dispatch
	appEventNames     string
	Version           string
	configData        []byte
}

func NewApp(config []byte) *App {
	this := &App{}
	this.configData = config
	var err error
	this.configFile, err = ini.Load(this.configData)
	if err != nil {
		fmt.Printf("configFile Fail to load %v", err)
	}
	this.Version = "2.0.0"
	this.dispatch = NewDispath(this)
	this.logManager = NewLogManager(this)
	this.logManager.Start()

	this.configManager = NewConfigManager(this, this.configData)
	this.moduleManager = NewModuleManager(this, this.configManager)
	this.errorManager = NewErrorManager(this)
	this.middlewareManager = NewMiddlewareManager(this)
	return this
}

func (this *App) Start() {
	go this.doWorker()
}

func (this *App) doWorker() {
	this.Debug("App", "load middlewares")
	func() {
		eventView := middlewares.NewEventView(this)
		this.middlewareManager.RegisterMiddleware(eventView)
		eventformat := middlewares.NewEventFormat(this)
		this.middlewareManager.RegisterMiddleware(eventformat)
		eventAuth := middlewares.NewEventView(this)
		this.middlewareManager.RegisterMiddleware(eventAuth)
	}()
	this.Debug("App", "Load modules")
	this.moduleManager.LoadModules()
	this.appEventNames = this.configManager.GlobalConfig.GetItem("subs")
	this.dispatch.SetEventNames(this.appEventNames)
	event := common_types.NewEvent("appReady", []byte{})
	this.Pub(event)
	go this.dispatch.Start()
}

func (this *App) Handel(event interfaces.Event) {
	//this.Debug("App", "App event"+event.GetEventName())
	switch event.GetEventName() {
	case "configChange":
		mconfig := &ModuleConfig{}
		json.Unmarshal(event.GetPayload(), mconfig)
		moduleName := mconfig.GetModuleName()
		oldModuleConfig := this.configManager.ModuleConfigs[moduleName]
		newConfigs := mconfig.GetItems()
		for key, val := range newConfigs {
			oldModuleConfig.SetItem(key, val)
		}
		break
	case "stopModule":
		moduleName := string(event.GetPayload())
		this.moduleManager.UnLoadModule(moduleName)
		break
	case "startModule":
		moduleName := string(event.GetPayload())
		moduleConfig := this.configManager.ModuleConfigs[moduleName]
		if moduleConfig == nil {
			moduleConfig = NewModuleConfig(moduleName, this.configManager.GlobalConfig)
		}
		this.moduleManager.LoadModule(moduleName, moduleConfig)
		break
	}
}

func (this *App) Sub(eventName string, module interfaces.Module) {
	if this.dispatch != nil {
		this.dispatch.Sub(eventName, module)
	} else {
		this.Error("App", "dispath unready")
	}
}

func (this *App) UnSub(eventName string, module interfaces.Module) {
	if this.dispatch != nil {
		this.dispatch.UnSub(eventName, module)
	} else {
		this.Error("App", "dispath unready")
	}
}

func (this *App) Pub(event interfaces.Event) {
	if this.middlewareManager.Handel(event) {
		this.dispatch.Pub(event)
	}
}

func (this *App) Info(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	this.logManager.Info(category, format, a...)
}

func (this *App) Warning(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	this.logManager.Waring(category, format, a...)
}

func (this *App) Error(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	this.logManager.Error(category, format, a...)
}

func (this *App) Debug(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	this.logManager.Debug(category, format, a...)
}

func (this *App) GetVersion() string {
	return this.Version
}

func (this *App) GetConfigItem(section, key string) (string, error) {
	sect, err := this.configFile.GetSection(section)
	if err != nil {
		return "", nil
	}
	key1, err := sect.GetKey(key)
	if err != nil {
		return "", err
	}
	return key1.String(), nil
}

func (this *App) GetDefaultConfigItem(key string) (string, error) {
	return this.GetConfigItem("DEFAULT", key)
}

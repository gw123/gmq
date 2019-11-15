package core

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

type App struct {
	errorManager      *ErrorManager
	moduleManager     *ModuleManager
	configManager     *ConfigManager
	middlewareManager *MiddlewareManager
	logManager        *LogManager
	dispatch          *Dispatch
	appEventNames     string
	Version           string
	configData        *viper.Viper
	DbPool            interfaces.DbPool
	RedisPool         interfaces.RedisPool
	CacheManager      interfaces.CacheManager
	Services          map[string]interfaces.Service
}

func NewApp(viper2 *viper.Viper) *App {
	this := &App{}
	this.configData = viper2

	this.Version = "1.0.0"
	this.dispatch = NewDispath(this)
	this.logManager = NewLogManager(this)
	this.logManager.Start()
	this.Services = make(map[string]interfaces.Service)
	this.configManager = NewConfigManager(this, viper2)
	this.moduleManager = NewModuleManager(this, this.configManager)
	this.errorManager = NewErrorManager(this)
	this.middlewareManager = NewMiddlewareManager(this)

	this.DbPool = NewDbPool()
	this.RedisPool = NewRedisPool(this.configData.GetStringMap("redisPool"), this)
	this.LoadDb()

	this.CacheManager = NewCacheManager(this)
	return this
}

func (app *App) Start() {
	go app.doWorker()
}

//加载数据库配置
func (app *App) LoadDb() {
	configs := app.configData.GetStringMap("dbpool")
	reg := regexp.MustCompile(`^\$\{(.*)\}$`)
	for key, config := range configs {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			if key != "default" {
				app.Info("App", "数据库配置文件 格式不支持")
			}
			continue
		}
		drive, ok := configMap["drive"].(string)
		if !ok {
			drive = "mysql"
		}
		if reg.MatchString(drive) {
			arrs := reg.FindStringSubmatch(drive)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				drive = os.Getenv(arrs[1])
			}
		}
		host, ok := configMap["host"].(string)
		if !ok {
			host = "127.0.0.1"
		}
		if reg.MatchString(host) {
			arrs := reg.FindStringSubmatch(host)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				host = os.Getenv(arrs[1])
			}
		}

		port, ok := configMap["port"].(string)
		if !ok || port == "" {
			port = "3306"
		}
		if reg.MatchString(port) {
			arrs := reg.FindStringSubmatch(port)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				port = os.Getenv(arrs[1])
			}
		}

		database, ok := configMap["database"].(string)
		if !ok {
			database = ""
		}
		if reg.MatchString(database) {
			arrs := reg.FindStringSubmatch(database)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				database = os.Getenv(arrs[1])
			}
		}

		username, ok := configMap["username"].(string)
		if !ok {
			username = "root"
		}
		if reg.MatchString(username) {
			arrs := reg.FindStringSubmatch(username)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				username = os.Getenv(arrs[1])
			}
		}

		password, ok := configMap["password"].(string)
		if !ok {
			password = ""
		}
		if reg.MatchString(password) {
			arrs := reg.FindStringSubmatch(password)
			if len(arrs) > 1 {
				app.Debug("app", "读取环境变量 %s", arrs[1])
				password = os.Getenv(arrs[1])
			}
		}

		app.Info("App", "load db name:%s ,database:%s", key, database)
		db, err := app.DbPool.NewDb(
			drive,
			host,
			port,
			database,
			username,
			password);
		if err != nil {
			app.Warn("App", "db load error, %s: ", err.Error())
		}

		maxIdles := app.configData.GetInt("dbpool." + key + ".max_idles")
		if maxIdles == 0 {
			maxIdles = 3
		}
		db.DB().SetMaxIdleConns(maxIdles)

		maxOpens := app.configData.GetInt("dbpool." + key + ".max_opens")
		if maxOpens == 0 {
			maxOpens = 30
		}
		app.Debug("Db", "set maxOpens %d,set maxIdles %d", maxOpens, maxIdles)
		db.DB().SetMaxIdleConns(maxOpens)

		app.DbPool.SetDb(key, db)
	}

	//set default db
	defaultDBkey, ok := configs["default"].(string)
	if !ok {
		return
	}

	app.Debug("App", "default DB: %s", defaultDBkey)

	db, err := app.DbPool.GetDb(defaultDBkey)
	if err != nil {
		app.Warn("App", "cant found db default config :%s", err.Error())
	} else {
		app.DbPool.SetDb("default", db)
	}
}

func (app *App) doWorker() {
	app.Debug("App", "Load modules")
	app.moduleManager.LoadModules()
	app.appEventNames = "stopModule,startModule,configChange"
	app.dispatch.SetEventNames(app.appEventNames)
	event := gmsg.NewEvent("appReady", []byte{})
	app.Pub(event)
	go app.dispatch.Start()
}

func (app *App) Handel(event interfaces.Msg) {
	//app.Debug("App", "App event"+event.GetEventName())
	switch event.GetEventName() {
	case "configChange":
		mconfig := &ModuleConfig{}
		json.Unmarshal(event.GetPayload(), mconfig)
		moduleName := mconfig.GetModuleName()
		oldModuleConfig := app.configManager.ModuleConfigs[moduleName]
		newConfigs := mconfig.GetItems()
		for key, val := range newConfigs {
			oldModuleConfig.SetItem(key, val)
		}
		break
	case "stopModule":
		moduleName := string(event.GetPayload())
		app.moduleManager.UnLoadModule(moduleName)
		break
	case "startModule":
		moduleName := string(event.GetPayload())
		moduleConfig := app.configManager.ModuleConfigs[moduleName]
		if moduleConfig == nil {
			moduleConfig = NewModuleConfig(moduleName, app.configManager.GlobalConfig)
		}
		app.moduleManager.LoadModule(moduleName, moduleConfig)
		break
	}
}

func (app *App) Sub(eventName string, module interfaces.Module) {
	if app.dispatch != nil {
		app.dispatch.Sub(eventName, module)
	} else {
		app.Error("App", "dispath unready")
	}
}

func (app *App) UnSub(eventName string, module interfaces.Module) {
	if app.dispatch != nil {
		app.dispatch.UnSub(eventName, module)
	} else {
		app.Error("App", "dispath unready")
	}
}

func (app *App) Pub(event interfaces.Msg) {
	if app.middlewareManager.Handel(event) {
		app.dispatch.Pub(event)
	}
}

func (app *App) Info(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	app.logManager.Info(category, format, a...)
}

func (app *App) Warn(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	app.logManager.Warn(category, format, a...)
}

func (app *App) Error(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	app.logManager.Error(category, format, a...)
}

func (app *App) Debug(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	app.logManager.Debug(category, format, a...)
}

func (app *App) GetVersion() string {
	return app.Version
}

func (app *App) GetConfigItem(section, key string) (string, error) {
	sect := app.configData.GetStringMapString(section)
	if sect != nil {
		return "", nil
	}
	val, ok := sect[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (app *App) GetConfig() *viper.Viper {
	return app.configData
}

func (app *App) GetAppConfigItem(key string) (string, error) {
	return app.GetConfigItem("app", key)
}

/***/
func (app *App) LoadModuleProvider(provider interfaces.ModuleProvider) {
	app.moduleManager.LoadModuleProvider(provider)
}

//获取数据库信息
func (app *App) GetDb(dbname string) (*gorm.DB, error) {
	if dbname == "" {
		return app.GetDefaultDb()
	}
	return app.DbPool.GetDb(dbname)
}

//获取默认数据库
func (app *App) GetDefaultDb() (*gorm.DB, error) {
	dbname := "default"
	return app.GetDb(dbname)
}

func (app *App) RegisterService(name string, service interfaces.Service) {
	if name == "" {
		name = service.GetServiceName()
	}
	if service == nil {
		app.Warn("App", "service is nil")
		return
	}
	app.Services[name] = service
	return
}

func (app *App) GetService(name string) interfaces.Service {
	return app.Services[name]
}

func (app *App) GetRedis(name string) (*redis.Client, error) {
	return app.RedisPool.GetDb(name)
}

func (app *App) GetDefaultRedis() (*redis.Client, error) {
	return app.RedisPool.GetDb("default")
}

func (app *App) GetCacheManager() (interfaces.CacheManager, error) {
	if app.CacheManager == nil {
		return nil, errors.New("app.cacheManager is nil")
	}
	return app.CacheManager, nil
}

func (app *App) GetLogger() (interfaces.Logger) {
	return app.logManager
}

func (app *App) Write(buf []byte) (int, error) {
	return app.logManager.Write(buf)
}

package core

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/common_types"
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
	return this
}

func (this *App) Start() {
	go this.doWorker()
}

//加载数据库配置
func (this *App) LoadDb() {
	configs := this.configData.GetStringMap("dbpool")
	reg := regexp.MustCompile(`^\$\{(.*)\}$`)
	for key, config := range configs {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			if key != "default" {
				this.Info("App", "数据库配置文件 格式不支持")
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
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
				this.Debug("app", "读取环境变量 %s", arrs[1])
				password = os.Getenv(arrs[1])
			}
		}

		this.Info("App", "load db name:%s ,database:%s", key, database)
		db, err := this.DbPool.NewDb(
			drive,
			host,
			port,
			database,
			username,
			password);
		if err != nil {
			this.Warn("App", "db load error, %s: ", err.Error())
		}

		maxIdles := this.configData.GetInt("dbpool." + key + ".max_idles")
		if maxIdles == 0 {
			maxIdles = 3
		}
		db.DB().SetMaxIdleConns(maxIdles)

		maxOpens := this.configData.GetInt("dbpool." + key + ".max_opens")
		if maxOpens == 0 {
			maxOpens = 30
		}
		this.Debug("Db", "set maxOpens %d,set maxIdles %d", maxOpens,maxIdles)
		db.DB().SetMaxIdleConns(maxOpens)

		this.DbPool.SetDb(key, db)
	}

	//set default db
	defaultDBkey, ok := configs["default"].(string)
	if !ok {
		return
	}

	this.Debug("App", "default DB: %s", defaultDBkey)

	db, err := this.DbPool.GetDb(defaultDBkey)
	if err != nil {
		this.Warn("App", "cant found db default config :%s", err.Error())
	} else {
		this.DbPool.SetDb("default", db)
	}
}

func (this *App) doWorker() {
	this.Debug("App", "Load modules")
	this.moduleManager.LoadModules()
	this.appEventNames = "stopModule,startModule,configChange"
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

func (this *App) Warn(category string, format string, a ...interface{}) {
	if format == "" {
		return
	}
	this.logManager.Warn(category, format, a...)
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
	sect := this.configData.GetStringMapString(section)
	if sect != nil {
		return "", nil
	}
	val, ok := sect[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (this *App) GetConfig() *viper.Viper {
	return this.configData
}

func (this *App) GetAppConfigItem(key string) (string, error) {
	return this.GetConfigItem("app", key)
}

/***/
func (this *App) LoadModuleProvider(provider interfaces.ModuleProvider) {
	this.moduleManager.LoadModuleProvider(provider)
}

//获取数据库信息
func (this *App) GetDb(dbname string) (*gorm.DB, error) {
	if dbname == "" {
		return this.GetDefaultDb()
	}
	return this.DbPool.GetDb(dbname)
}

//获取默认数据库
func (this *App) GetDefaultDb() (*gorm.DB, error) {
	dbname := "default"
	return this.GetDb(dbname)
}

func (this *App) RegisterService(name string, service interfaces.Service) {
	if name == "" {
		name = service.GetServiceName()
	}
	if service == nil {
		this.Warn("App", "service is nil")
		return
	}
	this.Services[name] = service
	return
}

func (this *App) GetService(name string) interfaces.Service {
	return this.Services[name]
}

func (this *App) GetRedis(name string) (*redis.Client, error) {
	return this.RedisPool.GetDb(name)
}

func (this *App) GetDefaultRedis() (*redis.Client, error) {
	return this.RedisPool.GetDb("default")
}

func (this *App) GetLogger() (interfaces.Logger) {
	return this.logManager
}

func (this *App) Write(buf []byte) (int, error) {
	return this.logManager.Write(buf)
}

package gmq

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type App interface {
	Logger
	//获取版本号
	GetVersion() string
	//发布消息
	Pub(event Msg)
	//订阅消息
	Sub(eventName string, module Module)
	//取消订阅
	UnSub(eventName string, module Module)
	//获取全局app配置信息
	GetAppConfigItem(key string) (val string, err error)
	//获取Viper对象
	GetConfig() *viper.Viper
	//处理消息
	Handel(event Msg)
	//加载模块提供者
	LoadModuleProvider(provider ModuleProvider)

	//获取数据库信息
	GetDb(dnname string) (*gorm.DB, error)
	//获取默认数据库
	GetDefaultDb() (*gorm.DB, error)
	//获取reids
	GetRedis(dnname string) (*redis.Client, error)
	//获取CacheManager
	GetCacheManager() (CacheManager, error)

	//获取默认reids
	GetDefaultRedis() (*redis.Client, error)
	RegisterService(name string, s Service)
	//RegisterServiceProvider(name string, s Service)
	GetService(name string) Service
	GetLogger() (Logger)
}

type Service interface {
	GetServiceName() string
}

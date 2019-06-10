package interfaces

import "github.com/jinzhu/gorm"

type App interface {
	Logger
	//获取版本号
	GetVersion() string
	//发布消息
	Pub(event Event)
	//订阅消息
	Sub(eventName string, module Module)
	//取消订阅
	UnSub(eventName string, module Module)
	//获取配置信息
	GetConfigItem(section, key string) (val string, err error)
	//获取全局配置信息
	GetDefaultConfigItem(key string) (val string, err error)
	//处理消息
	Handel(event Event)
	//加载模块提供者
	LoadModuleProvider(provider ModuleProvider)
	//获取数据库信息
	GetDb(dnname string) (*gorm.DB, error)
	//获取默认数据库
	GetDefaultDb() (*gorm.DB, error)
}

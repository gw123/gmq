package interfaces

const ModuleLoadSuccess = 0x01
const ModuleLoadFailed = 0x02
const ModuleUnLoad = 0x03
const ModuleUnInstall = 0x04

type Module interface {
	//订阅事件
	Init(app App, config ModuleConfig) error

	//取消事件订阅
	UnInit() error

	Push(event Event) error

	GetStatus() uint64

	GetModuleName() string

	//发布消息
	Pub(event Event)

	//订阅消息
	Sub(eventName string)

	//获取模块版本
	GetVersion() string

	GetEventNum() int
}

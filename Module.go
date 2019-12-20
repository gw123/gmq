package gmq

const ModuleLoadSuccess = 0x01
const ModuleLoadFailed = 0x02
const ModuleUnLoad = 0x03
const ModuleUnInstall = 0x04

type Module interface {
	GetApp() App

	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Debug(format string, a ...interface{})

	//订阅事件 初始化队列
	Init(app App, config ModuleConfig) error

	//启动前执行事件,模块可以在这里实现自己的方法
	BeforeStart() error

	//取消事件订阅
	UnInit() error

	Push(event Msg) error

	GetStatus() uint64

	GetModuleName() string

	//发布消息
	Pub(event Msg)

	//订阅消息, filter 过滤函数从接收到的消息中过滤不合法的消息
	Sub(eventName string, filter ...func(interface{}) bool)

	//获取模块版本
	GetVersion() string

	//开始处理事件
	Start()

	//处理事件
	Handle(event Msg) (err error)

	//定时调用方法
	Watch(index int)

	GetConfig() ModuleConfig
}

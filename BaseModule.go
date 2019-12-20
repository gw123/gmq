package gmq

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

const FullFlag_DropNew = 0x1
const FullFlag_DropOld = 0x0

type Handle func(e Msg) (err error)
type Watch func(index int)

type BaseModule struct {
	Config      ModuleConfig
	queue       chan Msg
	App         App
	isBusyMutex sync.Mutex
	muduleNmae  string
	fullFlag    int
	length      int
	status      int
	StopFlag    bool
	eventNames  []string
	Version     string
	module      Module
	Handle
	Watch
	cancelFun context.CancelFunc
	Ctx       context.Context
}

func (b *BaseModule) Init(app App, module Module, config ModuleConfig) error {
	b.App = app
	b.Config = config
	b.muduleNmae = config.GetModuleName()
	b.InitQueue(1024)
	b.StopFlag = false
	b.module = module
	b.Handle = module.Handle
	b.Watch = module.Watch
	evnetsStr := config.GetStringItem("subs")
	events := strings.Split(evnetsStr, ",")

	b.eventNames = events
	for _, eventName := range events {
		if eventName != "" {
			b.App.Sub(eventName, module)
		}
	}
	rootCtx := context.Background()
	ctx, cancelFun := context.WithCancel(rootCtx)
	b.cancelFun = cancelFun
	b.Ctx = ctx

	return nil
}

func (b *BaseModule) UnInit() error {
	b.StopFlag = true
	for _, eventName := range b.eventNames {
		b.App.UnSub(eventName, b.module)
	}
	b.Debug("BaseModule UnInit :" + b.GetModuleName())
	b.Stop()
	return nil
}

//这里不做处理留给子模块实现
func (b *BaseModule) BeforeStart() error {
	return nil
}

func (b *BaseModule) GetStatus() uint64 {
	return 1
}

func (b *BaseModule) GetVersion() string {
	return b.Version
}

func (b *BaseModule) InitQueue(length int) {
	b.length = length
	b.fullFlag = FullFlag_DropOld
	b.queue = make(chan Msg, length)
}

func (b *BaseModule) SetFullFlag(flag int) {
	b.fullFlag = flag
}

func (b *BaseModule) Push(event Msg) (err error) {
	b.isBusyMutex.Lock()
	defer b.isBusyMutex.Unlock()
	if len(b.queue) >= b.length {
		if b.fullFlag == FullFlag_DropOld {
			<-b.queue
		} else {
			return errors.New("queue is full")
		}
	}
	b.queue <- event
	return
}

func (b *BaseModule) startDaemon() {
	index := 0
	for ; ; {
		select {
		case _ = <-b.Ctx.Done():
			b.Info("Stop Module  start goroutine " + b.GetModuleName())
			return
		default:
			break
		}
		time.Sleep(time.Second)
		b.Watch(index)
		index++
		if index > 1000000 {
			index = 0
		}
	}
}

func (b *BaseModule) Stop() {
	if b.cancelFun != nil {
		b.Info("StopModule : " + b.GetModuleName())
		b.cancelFun()
		b.cancelFun = nil
	} else {
		b.Warning("StopModule : cancelFun not exist " + b.GetModuleName())
	}
}

func (b *BaseModule) Start() {
	go b.startDaemon()
	for ; ; {
		select {
		case _ = <-b.Ctx.Done():
			b.Info("StopModule : stop Start goroutine " + b.GetModuleName())
			return
			break

		case event := <-b.queue:
			if b.Handle != nil {
				err := b.Handle(event)
				if err != nil {
					b.Error("Handel 执行失败 " + event.GetEventName() + err.Error())
				}
			}
			break
		}
	}
}

func (b *BaseModule) GetModuleName() string {
	return b.muduleNmae
}

func (b *BaseModule) Info(content string, a ...interface{}) {
	b.App.Info(b.GetModuleName(), content, a ...)
}

func (b *BaseModule) Warning(content string, a ...interface{}) {
	b.App.Warn(b.GetModuleName(), content, a ...)
}

func (b *BaseModule) Error(content string, a ...interface{}) {
	b.App.Error(b.GetModuleName(), content, a ...)
}

func (b *BaseModule) Debug(content string, a ...interface{}) {
	b.App.Debug(b.GetModuleName(), content, a ...)
}

func (b *BaseModule) Pop() Msg {
	return <-b.queue
}

//发布消息
func (b *BaseModule) Pub(event Msg) {
	if event == nil {
		return
	}
	b.App.Pub(event)
}

//订阅消息
func (b *BaseModule) Sub(eventName string, filter ...func(interface{}) bool) {
	if eventName == "" {
		return
	}
	b.App.Sub(eventName, b.module)
}

//获取app对象
func (b *BaseModule) GetApp() App {
	return b.App
}

func (b *BaseModule) GetConfig() ModuleConfig {
	return b.Config
}

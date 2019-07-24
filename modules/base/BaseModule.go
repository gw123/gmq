package base

import (
	"errors"
	"sync"
	"strings"
	"context"
	"time"
	"github.com/gw123/GMQ/core/interfaces"
)

const FullFlag_DropNew = 0x1
const FullFlag_DropOld = 0x0

type Handle func(e interfaces.Event) (err error)
type Watch func(index int)

type BaseModule struct {
	Config      interfaces.ModuleConfig
	queue       chan interfaces.Event
	App         interfaces.App
	isBusyMutex sync.Mutex
	muduleNmae  string
	fullFlag    int
	length      int
	status      int
	StopFlag    bool
	eventNames  []string
	Version     string
	module      interfaces.Module
	Handle
	Watch
	cancelFun   context.CancelFunc
	Ctx         context.Context
}

func (this *BaseModule) Init(app interfaces.App, module interfaces.Module, config interfaces.ModuleConfig) error {
	this.App = app
	this.Config = config
	this.muduleNmae = config.GetModuleName()
	this.InitQueue(1024)
	this.StopFlag = false
	this.module = module
	this.Handle = module.Handle
	this.Watch = module.Watch
	evnetsStr := config.GetItem("subs")
	events := strings.Split(evnetsStr, ",")

	this.eventNames = events
	for _, eventName := range events {
		if eventName != "" {
			this.App.Sub(eventName, module)
		}
	}
	rootCtx := context.Background()
	ctx, cancelFun := context.WithCancel(rootCtx)
	this.cancelFun = cancelFun
	this.Ctx = ctx

	return nil
}

func (this *BaseModule) UnInit() error {
	this.StopFlag = true
	for _, eventName := range this.eventNames {
		this.App.UnSub(eventName, this.module)
	}
	this.Debug("BaseModule UnInit :" + this.GetModuleName())
	this.Stop()
	return nil
}

func (this *BaseModule) GetStatus() uint64 {
	return 1
}

func (this *BaseModule) GetVersion() string {
	return this.Version
}

func (this *BaseModule) InitQueue(length int) {
	this.length = length
	this.fullFlag = FullFlag_DropOld
	this.queue = make(chan interfaces.Event, length)
}

func (this *BaseModule) SetFullFlag(flag int) {
	this.fullFlag = flag
}

func (this *BaseModule) Push(event interfaces.Event) (err error) {
	this.isBusyMutex.Lock()
	defer this.isBusyMutex.Unlock()
	if len(this.queue) >= this.length {
		if this.fullFlag == FullFlag_DropOld {
			<-this.queue
		} else {
			return errors.New("queue is full")
		}
	}
	this.queue <- event
	return
}

func (this *BaseModule) startDaemon() {
	index := 0
	for ; ; {
		select {
		case _ = <-this.Ctx.Done():
			this.Info("Stop Module  start goroutine " + this.GetModuleName())
			return
		}

		time.Sleep(time.Millisecond * 100)
		this.Watch(index)
		index++
		if index > 1000000 {
			index = 0
		}
	}
}

func (this *BaseModule) Stop() {
	if this.cancelFun != nil {
		this.Info("StopModule : " + this.GetModuleName())
		this.cancelFun()
		this.cancelFun = nil
	} else {
		this.Warning("StopModule : cancelFun not exist " + this.GetModuleName())
	}
}

func (this *BaseModule) Start() {
	go this.startDaemon()
	for ; ; {
		select {
		case _ = <-this.Ctx.Done():
			this.Info("StopModule : stop Start goroutine " + this.GetModuleName())
			return
			break

		case event := <-this.queue:
			if this.Handle != nil {
				err := this.Handle(event)
				if err != nil {
					this.Error("Handel 执行失败 " + event.GetEventName() + err.Error())
				}
			}
			break
		}
	}
}

func (this *BaseModule) GetModuleName() string {
	return this.muduleNmae
}

func (this *BaseModule) Info(content string, a ...interface{}) {
	this.App.Info(this.GetModuleName(), content, a ...)
}

func (this *BaseModule) Warning(content string, a ...interface{}) {
	this.App.Warning(this.GetModuleName(), content, a ...)
}

func (this *BaseModule) Error(content string, a ...interface{}) {
	this.App.Error(this.GetModuleName(), content, a ...)
}

func (this *BaseModule) Debug(content string, a ...interface{}) {
	this.App.Debug(this.GetModuleName(), content, a ...)
}

func (this *BaseModule) Pop() interfaces.Event {
	return <-this.queue
}

//发布消息
func (this *BaseModule) Pub(event interfaces.Event) {
	if event == nil {
		return
	}
	this.App.Pub(event)
}

//订阅消息
func (this *BaseModule) Sub(eventName string) {
	if eventName == "" {
		return
	}
	this.App.Sub(eventName, this.module)
}

//获取app对象
func (this *BaseModule) GetApp() interfaces.App {
	return this.App
}

func (this *BaseModule) GetConfig() interfaces.ModuleConfig {
	return this.Config
}

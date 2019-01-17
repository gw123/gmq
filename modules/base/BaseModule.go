package base

import (
	"github.com/gw123/GMQ/core/interfaces"
	"errors"
	"sync"
	"time"
	"strings"
)

const FullFlag_DropNew = 0x1
const FullFlag_DropOld = 0x0

type BaseModule struct {
	Config      interfaces.ModuleConfig
	queue       chan interfaces.Event
	App         interfaces.App
	isBusyMutex sync.Mutex
	signal      chan int
	muduleNmae  string
	fullFlag    int
	length      int
	status      int
	StopFlag    bool
	eventNames  []string
	Version     string
}

func (this *BaseModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.App = app
	this.Config = config
	this.muduleNmae = config.GetModuleName()
	this.InitQueue(1024)
	this.StopFlag = false
	this.signal = make(chan int, 1)
	evnetsStr := config.GetItem("subs")
	events := strings.Split(evnetsStr, ",")
	this.eventNames = events
	for _, eventName := range events {
		this.App.Sub(eventName, this)
	}
	return nil
}

func (this *BaseModule) UnInit() error {
	this.StopFlag = true
	for _, eventName := range this.eventNames {
		this.App.UnSub(eventName, this)
	}
	return nil
}

func (this *BaseModule) GetStatus() uint64 {
	this.isBusyMutex.Lock()
	defer this.isBusyMutex.Unlock()
	return 1
}

func (this *BaseModule) GetVersion() string {
	return this.Version
}

func (this *BaseModule) GetEventNum() int {
	this.isBusyMutex.Lock()
	defer this.isBusyMutex.Unlock()
	return len(this.queue)
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

func (this *BaseModule) Handel(event interfaces.Event) (err error) {
	return
}

func (this *BaseModule) Watch() (event interfaces.Event) {
	return
}

func (this *BaseModule) StartDaemon() {
	go func() {
		for ; ; {
			if this.StopFlag {
				break
			}
			event := this.Watch()
			if event != nil {
				this.App.Pub(event)
			}
			time.Sleep(time.Second)
		}
	}()
}

func (this *BaseModule) Start() {
	go func() {
		for ; ; {
			if this.StopFlag {
				break
			}
			event := this.Pop()
			err := this.Handel(event)

			if err != nil {
				this.Error(event.GetMsgId() + " " + event.GetEventName() + " 执行失败 " + err.Error())
			} else {
				this.Info(event.GetMsgId() + " " + event.GetEventName() + " 执行成功")
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func (this *BaseModule) Pop() interfaces.Event {
	//this.isBusyMutex.Lock()
	//defer this.isBusyMutex.Unlock()
	return <-this.queue
}

func (this *BaseModule) GetModuleName() string {
	return this.muduleNmae
}

func (this *BaseModule) Info(content string) {
	this.App.Info(this.GetModuleName(), content)
}

func (this *BaseModule) Warning(content string) {
	this.App.Warning(this.GetModuleName(), content)
}

func (this *BaseModule) Error(content string) {
	this.App.Error(this.GetModuleName(), content)
}

func (this *BaseModule) Debug(content string) {
	this.App.Debug(this.GetModuleName(), content)
}

//发布消息
func (this *BaseModule) Pub(event interfaces.Event) {
	if event == nil {
		return
	}
	event.SetSourceModule(this.muduleNmae)
	this.App.Pub(event)
}

//订阅消息
func (this *BaseModule) Sub(eventName string) {
	if eventName == "" {
		return
	}
	this.App.Sub(eventName, this)
}

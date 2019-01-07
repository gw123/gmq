package core

import (
	"github.com/gw123/GMQ/core/interfaces"
	"sync"
	"strings"
	"time"
)

type Dispatch struct {
	EventQueueBinds map[string][]interfaces.Module
	EventQueues     interfaces.EventQueue
	app             interfaces.App
	mutex           sync.Mutex
	appEventNames   []string
}

func NewDispath(app interfaces.App) *Dispatch {
	this := new(Dispatch)
	this.app = app
	this.EventQueueBinds = make(map[string][]interfaces.Module)
	this.EventQueues = NewEventQueue(app)
	return this
}

func (this *Dispatch) SetEventNames(eventNames string) {
	this.appEventNames = strings.Split(eventNames, ",")
}

func (this *Dispatch) Start() {
	go func() {
		for ; ; {
			time.Sleep(time.Millisecond)
			event, err := this.EventQueues.Pop()
			//if event != nil {
			//	this.app.Debug("Dispath", "Start pop:"+event.GetMsgId()+" : "+event.GetEventName())
			//}
			if err != nil {
				if err.Error() == "队列为空" {
				} else {
					this.app.Warning("app", "出队异常:"+err.Error())
				}
				continue
			}

			if event == nil {
				this.app.Warning("app", "出队异常:event为nil")
				continue
			}

			eventName := event.GetEventName()
			modules := this.EventQueueBinds[eventName]
			//this.app.Debug("Dispatch",fmt.Sprintf("Bingdings modules len %d", len(modules)))
			for _, module := range modules {
				//this.app.Debug("Dispatch",fmt.Sprintf("Bingding moduleName %s", module.GetModuleName()))
				err := module.Push(event)
				if err != nil {
					this.app.Warning(module.GetModuleName(), "模块队列异常Push失败"+err.Error())
				}
			}
		}
	}()
}

func (this *Dispatch) PushToModule(event interfaces.Event) {
	eventName := event.GetEventName()
	modules := this.EventQueueBinds[eventName]
	//this.app.Debug("Dispatch",fmt.Sprintf("Bingdings modules len %d", len(modules)))
	for _, module := range modules {
		//this.app.Debug("Dispatch",fmt.Sprintf("Bingding moduleName %s", module.GetModuleName()))
		err := module.Push(event)
		if err != nil {
			this.app.Warning(module.GetModuleName(), "模块队列异常Push失败"+err.Error())
		}
	}
}

func (this *Dispatch) handel(event interfaces.Event) {
	app, ok := this.app.(*App)
	if ok {
		app.Handel(event)
	}
}

func (this *Dispatch) Sub(eventName string, module interfaces.Module) {
	if eventName == "" || eventName == " " {
		this.app.Warning("Dispatch", "Sub eventName 为空")
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	modules := this.EventQueueBinds[eventName]
	for _, m := range modules {
		if m.GetModuleName() == module.GetModuleName() {
			this.app.Warning("sub", m.GetModuleName()+"已经订阅"+eventName)
			return
		}
	}
	this.EventQueueBinds[eventName] = append(this.EventQueueBinds[eventName], module)
}

func (this *Dispatch) UnSub(eventName string, module interfaces.Module) {
	if eventName == "" {
		this.app.Warning("Dispatch", "UnSub eventName 为空")
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	modules := this.EventQueueBinds[eventName]
	for index, m := range modules {
		if m.GetModuleName() == module.GetModuleName() {
			modules[index] = nil
			this.app.Info("sub", m.GetModuleName()+"取消订阅"+eventName)
			return
		}
	}
}

func (this *Dispatch) Pub(event interfaces.Event) {
	if event.GetEventName() == "" {
		this.app.Warning("Dispatch", "Pub eventName 为空")
		return
	}
	//this.app.Debug("Dispath", event.GetMsgId()+" : "+event.GetEventName())
	//this.EventQueues.Push(event)
	this.PushToModule(event)
}

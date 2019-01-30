package debugModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/common/common_types"
	"fmt"
	"time"
	"strconv"
	"github.com/fpay/erp-client-s/common"
)

type DebugModule struct {
	base.BaseModule
	status bool
	count  int
}

func NewDebugModule() *DebugModule {
	this := new(DebugModule)
	return this
}

//订阅 mqtt_upload  mqtt_log
func (this *DebugModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	app.Sub("debug", this)
	//config.GetGlobalItem("")
	go this.StartDaemon()
	return nil
}

func (this *DebugModule) UnInit() error {
	this.App.UnSub("debug", this)
	this.StopFlag = true
	return nil
}

func (this *DebugModule) GetStatus() uint64 {
	if this.status {
		return 1
	} else {
		return 0
	}
}

func (this *DebugModule) Watch() (event interfaces.Event) {
	this.count++
	if this.count%20 == 0 {
		event = common.NewEvent("testlib", []byte("10"+strconv.Itoa(this.count)))
		this.App.Pub(event)
	}
	if this.count%40 == 0 {
		event = common_types.NewEvent("printe", []byte("test printe"))
		this.App.Pub(event)
	}

	if this.count%100 == 10 {
		//event = common.NewEvent("log", []byte("log event"))
		//this.App.Pub(event)
		this.Error("测试上报错误")
	}
	return event
}

func (this *DebugModule) Handel(event interfaces.Event) error {
	str := fmt.Sprintf("msgId:%s,eventName:%s,payload:%s", event.GetMsgId(), event.GetEventName(), event.GetPayload())
	this.Debug(str)
	return nil
}

func (this *DebugModule) StartDaemon() {
	for ; ; {
		if this.StopFlag {
			break
		}
		event := this.Watch()
		if event != nil {
			//新的事件
			this.App.Pub(event)
		}
		time.Sleep(time.Second)
	}
}

func (this *DebugModule) Start() {
	for ; ; {
		if this.StopFlag {
			break
		}
		event := this.Pop()
		err := this.Handel(event)

		if err != nil {
			//执行失败
			//replay := common.NewResultEvent([]byte("失败" + err.Error()))
			//this.App.Pub(replay)
		} else {
			//执行成功
			//replay := common.NewResultEvent([]byte("成功"))
			//this.App.Pub(replay)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

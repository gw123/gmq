package testModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"fmt"
	"time"
)

type testModule struct {
	base.BaseModule
}

func NewtestModule() *testModule {
	this := new(testModule)
	return this
}

func (this *testModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	app.Sub("printe", this)
	go this.Start()
	return nil
}

func (this *testModule) UnInit() error {
	this.App.UnSub("printe", this)
	return nil
}

func (this *testModule) GetStatus() uint64 {
	return 1
}

func (this *testModule) Start() {
	go func() {
		//this.App.Pub()
	}()
	for ; ; {
		event := this.BaseModule.Pop()
		err := this.service(event)
		if err != nil {
			//执行失败
			fmt.Println("testModule service " + err.Error())
			//replay := NewPrinterResultEvent(event.GetMsgId(), "打印失败"+err.Error())
			//this.App.Pub(replay)
		} else {
			//执行成功
			//replay := NewResultEvent(event.GetMsgId(), "打印成功")
			//this.App.Pub(replay)
		}
		time.Sleep(time.Second)
	}
}

func (this *testModule) service(event interfaces.Event) error {
	this.Info(event.GetEventName() + ", " + event.GetMsgId() + " ," + string(event.GetPayload()))
	return nil
}

package mnsModule

import (
	"github.com/fpay/erp-client-s/modules/base"
	"github.com/fpay/erp-client-s/interfaces"
)

type MqttModule struct {
	base.BaseModule
	IsConnectAliIot bool
}

func NewMqttModule() *MqttModule {
	this := new(MqttModule)
	return this
}

//订阅 mqtt_upload  mqtt_log
func (this *MqttModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	//this.SetApp(app)
	this.BaseModule.Init(app, config)
	app.Sub("reply", this)
	app.Sub("sendMqttMsg", this)
	return nil
}

func (this *MqttModule) UnInit() error {
	this.App.UnSub("reply", this)
	this.App.UnSub("sendMqttMsg", this)
	return nil
}

func (this *MqttModule) Push(event interfaces.Event) (err error) {
	err = this.BaseModule.Push(event)
	return
}

func (this *MqttModule) GetStatus() uint64 {
	if this.IsConnectAliIot {
		return 1
	} else {
		return 0
	}
}

func (this *MqttModule) service(event interfaces.Event) error {
	this.Debug("service " + event.GetEventName() + " " + event.GetSourceModule() + " " + event.GetMsgId() + " " + string(event.GetPayload()))
	return nil
}

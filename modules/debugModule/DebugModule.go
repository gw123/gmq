package debugModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"fmt"
)

type DebugModule struct {
	base.BaseModule
}

func NewDebugModule() *DebugModule {
	this := new(DebugModule)
	return this
}

//订阅 mqtt_upload  mqtt_log
func (this *DebugModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	app.Sub("debug", this)
	return nil
}

func (this *DebugModule) Handle(event interfaces.Event) error {
	return nil
}

func (this *DebugModule) Watch(index int) {

	return
}

func (this *DebugModule) Handel(event interfaces.Event) error {
	str := fmt.Sprintf("msgId:%s,eventName:%s,payload:%s", event.GetMsgId(), event.GetEventName(), event.GetPayload())
	this.Debug(str)
	return nil
}



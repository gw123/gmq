package serverNodeModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
)

type ServerNodeModule struct {
	base.BaseModule
	eventNames []string
	authToken  string
	clientName string
	serverHost string
	nodeName   string
}

func NewServerNodeModule() *ServerNodeModule {
	this := new(ServerNodeModule)
	return this
}

func (this *ServerNodeModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	this.serverHost = this.Config.GetStringItem("serverHost")
	this.nodeName = this.Config.GetStringItem("nodeName")
	//this.addr = this.Config.GetStringItem("nodeName")
	return nil
}

func (this *ServerNodeModule) Handle(event interfaces.Event) error {
	return nil
}

func (this *ServerNodeModule) Watch(count int) {
	if (count%10 == 0) {

	}
	return
}

func (this *ServerNodeModule) Start() {
	//初始化数据库

}

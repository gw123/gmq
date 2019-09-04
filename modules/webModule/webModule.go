package webModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/modules/webModule/server"
)

type WebModule struct {
	base.BaseModule
	eventNames []string
	authToken  string
	clientName string
	port       int
	addr       string
	controller *controllers.WsController
	server     *server.Server
}

func NewWebModule() *WebModule {
	this := new(WebModule)
	return this
}

func (this *WebModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	this.port = this.Config.GetIntItem("port")
	this.addr = this.Config.GetItem("addr")
	return nil
}

func (this *WebModule) Handle(event interfaces.Event) error {
	return nil
}

func (this *WebModule) Watch(count int) {
	return
}

func (this *WebModule) Start() {
	//初始化数据库
	err := autoMigrate(this.GetApp())
	if err != nil {
		this.Error("autoMigrate error %s .", err.Error())
	}

	port := this.Config.GetIntItem("port")
	addr := this.Config.GetItem("addr")
	this.server = server.NewServer(addr, port, this)
	go this.server.Start()
	this.BaseModule.Start()
}

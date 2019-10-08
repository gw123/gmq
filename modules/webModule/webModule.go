package webModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/modules/webModule/server"
	"github.com/pkg/errors"
)

type WebModule struct {
	base.BaseModule
	eventNames []string
	authToken  string
	clientName string
	controller *controllers.WsController
	server     *server.Server
	addr       string
}

func NewWebModule() *WebModule {
	this := new(WebModule)
	return this
}

func (this *WebModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	this.addr = this.Config.GetStringItem("bindAddr")
	if this.addr == "" {
		return errors.New("加载失败,bindAddr")
	}
	return nil
}

func (this *WebModule) Handle(event interfaces.Event) error {
	return nil
}

func (this *WebModule) Watch(count int) {
	return
}

func (this *WebModule) Start() {
	this.server = server.NewServer(this.addr, this)
	go this.server.Start()
	this.BaseModule.Start()
}

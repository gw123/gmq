package webModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/labstack/echo"
	"fmt"
	"github.com/gw123/GMQ/common/common_types"
	"strconv"
)

type WebModule struct {
	base.BaseModule
	eventNames []string
	authToken  string
	clientName string
	port       int
	addr       string
	controller *controllers.IndexController
	echo       *echo.Echo
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
	if count%20 == 0 {
		event := common_types.NewEvent("testlib", []byte("10"+strconv.Itoa(count)))
		this.App.Pub(event)
	}

	if count%40 == 0 {
		event := common_types.NewEvent("printe", []byte("test printe"))
		this.App.Pub(event)
	}

	if count%100 == 10 {
		this.Error("测试上报错误")
	}
	return
}

func (this *WebModule) Start() {
	go this.InitWebServer()
	this.BaseModule.Start()
}

func (this *WebModule) InitWebServer() error {
	controller := controllers.NewIndexController(this)
	this.controller = controller
	e := echo.New()
	e.GET("/message", controller.Message)
	e.GET("/sendMessage", controller.SendMessage)
	//e.GET("/", controller.)
	addr := fmt.Sprintf("%s:%d", this.addr, this.port)
	this.Info("端口监听在:  %s", addr)
	this.echo = e
	e.Start(addr)
	return nil
}

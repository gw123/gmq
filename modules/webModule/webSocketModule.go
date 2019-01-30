package webModule

import (
	"github.com/fpay/erp-client-s/modules/base"
	"github.com/fpay/erp-client-s/interfaces"
	"time"
	"strconv"
	"github.com/fpay/erp-client-s/modules/webModule/controllers"
	"github.com/labstack/echo"
	"fmt"
)

type WebModule struct {
	base.BaseModule
	eventNames []string
	authToken  string
	clientName string
	port       int
	controller *controllers.IndexController
	echo       *echo.Echo
}

func NewWebSocketModule() *WebModule {
	this := new(WebModule)
	return this
}

func (this *WebModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	var err error
	this.port, err = strconv.Atoi(this.Config.GetGlobalItem("port"))
	if err != nil {
		this.port = 17335
	}
	this.Debug("webModule subs: "+this.Config.GetItem("subs"))
	return nil
}

func (this *WebModule) UnInit() error {
	this.BaseModule.UnInit()
	return nil
}

func (this *WebModule) GetStatus() uint64 {
	return 1
}

func (this *WebModule) Start() {
	go this.InitWebSocket()
	time.Sleep(time.Second)
	for ; ; {
		event := this.BaseModule.Pop()
		err := this.service(event)
		if err != nil {
			this.Warning("WebModule service " + err.Error())
		}
		time.Sleep(time.Millisecond)
	}
}

func (this *WebModule) service(event interfaces.Event) error {
	this.Debug(event.GetEventName() + ", " + event.GetMsgId() + " ," + string(event.GetPayload()))
	this.controller.SendMessage(event)
	return nil
}

func (this *WebModule) InitWebSocket() error {
	controller := controllers.NewIndexController(this)
	this.controller = controller
	e := echo.New()
	e.GET("/message", controller.Message)
	e.GET("/", controller.Index)
	addr := fmt.Sprintf("127.0.0.1:%d", this.port)
	this.Info("端口监听在: " + addr)
	this.echo = e
	e.Start(addr)
	return nil
}

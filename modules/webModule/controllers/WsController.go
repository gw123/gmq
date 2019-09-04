package controllers

import (
	"github.com/gw123/GMQ/modules/webModule/models"
	"github.com/gw123/GMQ/core/interfaces"
	"golang.org/x/net/websocket"
	"github.com/labstack/echo"
	context2 "golang.org/x/net/context"
	"sync"
	"github.com/gw123/GMQ/common/common_types"
	"time"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"net/http"
	"strconv"
	"fmt"
)

type WsController struct {
	BaseController
	WebSocketClientMap map[string]*models.WsClientModel
	Mutex              sync.Mutex
	webModule          interfaces.Module
}

func NewIndexController(module interfaces.Module) *WsController {
	temp := new(WsController)
	temp.webModule = module
	temp.WebSocketClientMap = make(map[string]*models.WsClientModel, 10)
	return temp
}

func (c *WsController) Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index", nil)
}

func (c *WsController) Ping(ctx echo.Context) error {
	sleep := ctx.QueryParam("sleep")
	sleepTime, err := strconv.Atoi(sleep)
	fmt.Println(sleep, sleepTime)
	if err == nil && sleepTime < 1000*120 {
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
	return c.Success(ctx, nil)
}

func (c *WsController) Status(ctx echo.Context) error {
	return c.Success(ctx, "success")
}

func (c *WsController) Message(ctx echo.Context) error {
	moduleName := ctx.QueryParam("moduleName")
	c.webModule.Info("New WsClient coming! moduleName:" + moduleName)
	if ctx.IsWebSocket() {
		websocket.Handler(func(ws *websocket.Conn) {
			//c.webModule.Debug("ws handel ->")
			client, ok := c.WebSocketClientMap[moduleName]
			if ok {
				stopEvent := common_types.NewEvent("stop", []byte("新的同名模块连接到来"))
				client.SendMsg(stopEvent)
				client.Stop()
			}
			context := context2.Background()
			client = models.NewWsClientModel(ws, context, c.webModule, moduleName)
			c.Mutex.Lock()
			c.WebSocketClientMap[moduleName] = client
			c.Mutex.Unlock()
			//c.webModule.Debug("ws handel ->2")
			client.Run()
		}).ServeHTTP(ctx.Response(), ctx.Request())
	} else {
		c.webModule.Info("Message: 非法请求")
		response := &webEvent.Response{
			Code: 1,
			Msg:  "非法请求",
		}
		ctx.JSON(500, response)
	}
	return nil
}

func (c *WsController) SendClientMessage(msg interfaces.Event) {
	var flag = false
	for !flag {
		for _, client := range c.WebSocketClientMap {
			if client == nil {
				continue
			}
			flag = true
			err := client.SendMsg(msg)
			if err != nil {
				c.webModule.Warning("snedMsg error " + err.Error())
			}
		}
		time.Sleep(time.Second)
	}
}

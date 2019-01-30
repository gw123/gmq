package webSocketModule

import (
	"github.com/gw123/GMQ/common/common_types"
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"golang.org/x/net/websocket"
	"qiniupkg.com/x/errors.v7"
	"encoding/json"
	"strings"
	"time"
	"sync"
	"fmt"
)

type WebSocketModule struct {
	eventNames   []string
	originUrl    string
	websocketUrl string
	authToken    string
	clientName   string
	Conn         *websocket.Conn
	base.BaseModule
}

func NewWebSocketModule() *WebSocketModule {
	this := new(WebSocketModule)
	return this
}

func (this *WebSocketModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	this.originUrl = config.GetItem("originUrl")
	this.websocketUrl = config.GetItem("websocketUrl")
	go this.InitWebSocket()
	return nil
}

func (this *WebSocketModule) UnInit() error {
	this.BaseModule.UnInit()
	this.Conn.Close()
	return nil
}

func (this *WebSocketModule) GetStatus() uint64 {
	return 1
}

func (this *WebSocketModule) Start() {
	go func() {
		//this.App.Pub()
	}()
	for ; ; {
		event := this.BaseModule.Pop()
		err := this.service(event)
		if err != nil {
			//执行失败
			fmt.Println("WebSocketModule service " + err.Error())
			//replay := NewPrinterResultEvent(event.GetMsgId(), "打印失败"+err.Error())
			//this.App.Pub(replay)
		} else {
			//执行成功
			//replay := NewPrinterResultEvent(event.GetMsgId(), "打印成功")
			//this.App.Pub(replay)
		}
		time.Sleep(time.Second)
	}
}

func (this *WebSocketModule) service(event interfaces.Event) error {
	this.Info(event.GetEventName() + ", " + event.GetMsgId() + " ," + string(event.GetPayload()))
	if this.Conn == nil {
		return errors.New("WebSocket 连接未建立")
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		this.Error("json.Marshal " + err.Error())
		return err
	}
	//	this.Debug("eventData:" + string(eventData))
	_, err = this.Conn.Write(eventData)
	if err != nil {
		this.Error("Conn.Write " + err.Error())
		return err
	}
	return nil
}

func (this *WebSocketModule) InitWebSocket() error {
	url := this.websocketUrl + "?authToken=token_gw123&clientName=innerLogServer"
	//this.Debug("url" + url)
	ws, err := websocket.Dial(url, "", this.originUrl)
	if err != nil {
		this.Error("InitWebSocket " + err.Error())
		return err
	}
	this.Conn = ws
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg sync.WaitGroup) {
		defer wg.Done()
		var msg = make([]byte, 1024*1024)
		for {
			length, err := ws.Read(msg)
			if err != nil {
				if strings.Contains(err.Error(), "forcibly closed") {
					break
				}
				this.Error("InitWebSocket " + err.Error())
				continue
			}
			event := &common_types.Event{}
			err = json.Unmarshal(msg[0:length], event)
			if err != nil {
				this.Error("InitWebSocket " + err.Error())
				continue

			}
			this.Push(event)
		}
	}(wg)
	wg.Wait()
	return nil
}

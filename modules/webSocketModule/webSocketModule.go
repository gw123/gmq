package webSocketModule

import (
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"golang.org/x/net/websocket"
	"encoding/json"
	"strings"
	"time"
	"sync"
	"fmt"
	"errors"
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

func (w *WebSocketModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	w.BaseModule.Init(app, w, config)
	w.originUrl = config.GetStringItem("originUrl")
	w.websocketUrl = config.GetStringItem("websocketUrl")
	go w.InitWebSocket()
	return nil
}

func (w *WebSocketModule) UnInit() error {
	w.BaseModule.UnInit()
	w.Conn.Close()
	return nil
}

func (w *WebSocketModule) GetStatus() uint64 {
	return 1
}

func (w *WebSocketModule) Start() {

	for ; ; {
		event := w.BaseModule.Pop()
		err := w.service(event)
		if err != nil {
			//执行失败
			fmt.Println("WebSocketModule service " + err.Error())
			//replay := NewPrinterResultEvent(event.GetMsgId(), "打印失败"+err.Error())
			//w.App.Pub(replay)
		} else {
			//执行成功
			//replay := NewPrinterResultEvent(event.GetMsgId(), "打印成功")
			//w.App.Pub(replay)
		}
		time.Sleep(time.Second)
	}
}

func (w *WebSocketModule) Handel(event interfaces.Msg) error {
	w.Info(event.GetEventName() + ", " + event.GetMsgId() + " ," + string(event.GetPayload()))
	if w.Conn == nil {
		return errors.New("WebSocket 连接未建立")
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		w.Error("json.Marshal " + err.Error())
		return err
	}
	//	w.Debug("eventData:" + string(eventData))
	_, err = w.Conn.Write(eventData)
	if err != nil {
		w.Error("Conn.Write " + err.Error())
		return err
	}
	return nil
}

func (w *WebSocketModule) InitWebSocket() error {
	url := w.websocketUrl + "?authToken=token_gw123&clientName=innerLogServer"
	//w.Debug("url" + url)
	ws, err := websocket.Dial(url, "", w.originUrl)
	if err != nil {
		w.Error("InitWebSocket " + err.Error())
		return err
	}
	w.Conn = ws
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
				w.Error("InitWebSocket " + err.Error())
				continue
			}
			event := &gmsg.Msg{}
			err = json.Unmarshal(msg[0:length], event)
			if err != nil {
				w.Error("InitWebSocket " + err.Error())
				continue

			}
			w.Push(event)
		}
	}(wg)
	wg.Wait()
	return nil
}

func (w *WebSocketModule) Handle(event interfaces.Msg) (err error) {
	panic("implement me")
}

func (w *WebSocketModule) Watch(index int) {
	panic("implement me")
}
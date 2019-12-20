package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"golang.org/x/net/websocket"
	"io"
	"math/rand"
	"sync"
	"time"
)



type WsClientModel struct {
	webSocket  *websocket.Conn
	Mutex      sync.Mutex
	Token      string
	ModuleName string
	context    context.Context
	webEvent   interfaces.Module
	runFlag    bool
}

func NewWsClientModel(conn *websocket.Conn, ctx context.Context, module interfaces.Module, ModuleName string) *WsClientModel {
	this := new(WsClientModel)
	this.webSocket = conn
	this.webEvent = module
	this.context = ctx
	//key := "ModuleName"
	this.ModuleName = ModuleName
	this.runFlag = true
	return this
}

func (ws *WsClientModel) IsSafe(event *webEvent.RequestEvent) bool {
	if event.Token == ws.Token {
		return true
	}
	return false
}

func (ws *WsClientModel) Stop() {
	ws.runFlag = false
	if ws.webSocket.IsClientConn() || ws.webSocket.IsServerConn() {
		_ = ws.webSocket.Close()
	}
}

func (ws *WsClientModel) Run() {
	for ws.runFlag {
		select {
		case <-ws.context.Done():
			ws.runFlag = false
			break
		default:
		}
		//ws.webEvent.Info("ws handel ->6")
		event, err := ws.ReadMsg()

		if err == io.EOF {
			ws.webEvent.Warning("ReadMsg: 连接断开")
			break
		}
		if err != nil {
			ws.webEvent.Warning("ReadMsg:" + err.Error())
			if ws.webSocket.IsClientConn() && ws.webSocket.IsServerConn() {
				continue
			} else {
				break
			}
		}
		ws.DealMsg(event)
	}
}

func (ws *WsClientModel) DealMsg(event *webEvent.RequestEvent) {
	ws.webEvent.Debug(fmt.Sprintf("EeventName:%s; ModuleName:%s; Payload:%s", event.EventName, ws.ModuleName, event.Payload))
	switch event.EventName {
	case "auth":
		ws.ModuleName = event.SourcModuleName
		rand.Seed(time.Now().UnixNano())
		ws.Token = fmt.Sprintf("%d", rand.Uint64())
		ev := webEvent.NewEvent("auth_reply", ws.Token)
		ws.SendMsg(ev)
		break
	default:
		if ws.IsSafe(event) {
			ws.webEvent.Pub(event)
		} else {
			ws.webEvent.Warning("未认证连接")
			ws.Stop()
		}
	}
}

func (ws *WsClientModel) ReadMsg() (*webEvent.RequestEvent, error) {
	var gbuffer = make([]byte, 1024*1024)
	n, err := ws.webSocket.Read(gbuffer)
	if err != nil {
		return nil, err
	}
	decodeBufer := gbuffer[:n]
	event := &webEvent.RequestEvent{}
	err = json.Unmarshal(decodeBufer, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (ws *WsClientModel) SendMsg(event interfaces.Msg) error {
	event2 := &webEvent.Event{}
	event2.EventName = event.GetEventName()
	event2.Payload = string(event.GetPayload())
	event2.MsgId = event.GetMsgId()
	eventData, err := json.Marshal(event2)
	if err != nil {
		return err
	}

	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	_, err = ws.webSocket.Write(eventData)
	if err != nil {
		return err
	}
	return nil
}

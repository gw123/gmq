package models

import (
	"golang.org/x/net/websocket"
	"sync"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"time"
	"fmt"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
)

var gbuffer = make([]byte, 1024*1024)

type WsClientModel struct {
	webSocket  *websocket.Conn
	Mutex      sync.Mutex
	Token      string
	ModuleName string
	context    context.Context
	webEvent  interfaces.Module
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

func (this *WsClientModel) IsSafe(event *webEvent.RequestEvent) bool {
	return true
	if event.Token == this.Token {
		return true
	}
	return false
}

func (this *WsClientModel) Stop() {
	this.runFlag = false
	if this.webSocket.IsClientConn() || this.webSocket.IsServerConn() {
		this.webSocket.Close()
	}
}

func (this *WsClientModel) Run() {
	for this.runFlag {
		select {
		case <-this.context.Done():
			this.runFlag = false
			break
		default:
		}
		//this.webEvent.Info("ws handel ->6")
		event, err := this.ReadMsg()

		if err == io.EOF {
			this.webEvent.Warning("ReadMsg: 连接断开")
			break
		}
		if err != nil {
			this.webEvent.Warning("ReadMsg:" + err.Error())
			if this.webSocket.IsClientConn() && this.webSocket.IsServerConn() {
				continue
			} else {
				break
			}
		}
		this.DealMsg(event)
	}
}

func (this *WsClientModel) DealMsg(event *webEvent.RequestEvent) {
	this.webEvent.Debug(fmt.Sprintf("EeventName:%s; ModuleName:%s; Payload:%s", event.EventName, this.ModuleName, event.Payload))
	switch event.EventName {
	case "auth":
		this.ModuleName = event.SourcModuleName
		rand.Seed(time.Now().UnixNano())
		this.Token = fmt.Sprintf("%d", rand.Uint64())
		ev := webEvent.NewEvent("auth_reply", this.Token)
		this.SendMsg(ev)
		break
	default:
		if this.IsSafe(event) {
			this.webEvent.Pub(event)
		} else {
			this.webEvent.Warning("未认证连接")
			this.Stop()
		}
	}
}

func (this *WsClientModel) ReadMsg() (*webEvent.RequestEvent, error) {
	n, err := this.webSocket.Read(gbuffer)
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

func (this *WsClientModel) SendMsg(event interfaces.Msg) error {
	event2 := &webEvent.Event{}
	event2.EventName = event.GetEventName()
	event2.Payload = string(event.GetPayload())
	event2.MsgId = event.GetMsgId()
	eventData, err := json.Marshal(event2)
	if err != nil {
		return err
	}

	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	_, err = this.webSocket.Write(eventData)
	if err != nil {
		return err
	}
	return nil
}

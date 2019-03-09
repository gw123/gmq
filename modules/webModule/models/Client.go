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
	"os"
	"github.com/gw123/GMQ/modules/webModule/common"
	"github.com/gw123/GMQ/core/interfaces"
)

var gbuffer = make([]byte, 1024*1024)

type WsClientModel struct {
	webSocket  *websocket.Conn
	Mutex      sync.Mutex
	Token      string
	ModuleName string
	context    context.Context
	webModule  interfaces.Module
	runFlag    bool
}

func NewWsClientModel(conn *websocket.Conn, ctx context.Context, module interfaces.Module, ModuleName string) *WsClientModel {
	this := new(WsClientModel)
	this.webSocket = conn
	this.webModule = module
	this.context = ctx
	//key := "ModuleName"
	this.ModuleName = ModuleName
	this.runFlag = true
	return this
}

func (this *WsClientModel) IsSafe(event *common.RequestEvent) bool {
	return true
	if event.Token == this.Token && event.ModuleName == this.ModuleName {
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
		//this.webModule.Info("ws handel ->6")
		event, err := this.ReadMsg()

		if err == io.EOF {
			this.webModule.Warning("ReadMsg: 连接断开")
			break
		}
		if err != nil {
			this.webModule.Warning("ReadMsg:" + err.Error())
			if this.webSocket.IsClientConn() && this.webSocket.IsServerConn() {
				continue
			} else {
				break
			}
		}
		this.DealMsg(event)
	}
}

func (this *WsClientModel) DealMsg(event *common.RequestEvent) {
	this.webModule.Debug(fmt.Sprintf("EeventName:%s; ModuleName:%s; Payload:%s", event.EventName, this.ModuleName, event.Payload))
	switch event.EventName {
	case "auth":
		this.ModuleName = event.ModuleName
		rand.Seed(time.Now().UnixNano())
		this.Token = fmt.Sprintf("%d", rand.Uint64())
		ev := common.NewEvent("auth_reply", this.Token)
		this.SendMsg(ev)
		break
	case "exit":
		this.webModule.Info("收到退出消息...")
		time.Sleep(time.Second * 2)
		os.Exit(0)
		break
	default:

		if this.IsSafe(event) {
			this.webModule.Pub(event)
		} else {
			this.webModule.Warning("未认证连接")
			this.Stop()
		}
	}
}

func (this *WsClientModel) ReadMsg() (*common.RequestEvent, error) {
	n, err := this.webSocket.Read(gbuffer)
	if err != nil {
		return nil, err
	}
	decodeBufer := gbuffer[:n]
	event := &common.RequestEvent{}
	err = json.Unmarshal(decodeBufer, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (this *WsClientModel) SendMsg(event interfaces.Event) error {
	event2 := &common.Event{}
	event2.EventName = event.GetEventName()
	event2.Payload = string(event.GetPayload())
	event2.MsgId = event.GetMsgId()
	event2.SetSourceModule(event.GetSourceModule())
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

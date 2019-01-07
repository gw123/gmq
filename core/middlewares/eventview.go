package middlewares

import (
	"github.com/gw123/GMQ/core/interfaces"
	"golang.org/x/net/websocket"
	"sync"
	"time"
	"fmt"
	"log"
)

type EventView struct {
	app interfaces.App
}

func NewEventView(app interfaces.App) *EventView {
	this := new(EventView)
	this.app = app
	return this
}

func (this *EventView) Handel(event interfaces.Event) bool {
	//this.app.Debug("eventView", event.GetMsgId()+":"+event.GetEventName())
	return true
}

func (this *EventView) GetAttachEventTypes() string {
	return "*"
}

func (this *EventView) SendLog() {
	var origin = "http://localhost:1323"
	var url = "ws://localhost:1323/sendLog"
	url = url + "?authToken=token_gw123&clientName=innerLogServer"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal("Connect", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg sync.WaitGroup) {
		defer wg.Done()
		for {
			time.Sleep(time.Second * 5)
			message := []byte("hello, world!你好\n")
			_, err = ws.Write(message)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Send: %s\n", message)
		}
	}(wg)

	wg.Add(1)
	go func(wg sync.WaitGroup) {
		defer wg.Done()
		for {
			var msg = make([]byte, 512)
			m, err := ws.Read(msg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Receive: %s\n", msg[:m])
		}
	}(wg)
	wg.Wait()
	ws.Close() //关闭连接
}
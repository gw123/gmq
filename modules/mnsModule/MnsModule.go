package mnsModule

import (
	"github.com/gogap/ali_mns"
	"github.com/gogap/logs"
	"strings"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
	"encoding/json"
	"fmt"
	"encoding/base64"
)

type MnsModule struct {
	base.BaseModule
	IsConnectAliIot bool
}

func NewMnsModule() *MnsModule {
	this := new(MnsModule)
	return this
}

//订阅 mqtt_upload  mqtt_log
func (this *MnsModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, config)
	this.Watch()
	return nil
}

func (this *MnsModule) UnInit() error {
	this.App.UnSub("sendMqttMsg", this)
	return nil
}

func (this *MnsModule) Push(event interfaces.Event) (err error) {
	err = this.BaseModule.Push(event)
	return
}

func (this *MnsModule) GetStatus() uint64 {
	if this.IsConnectAliIot {
		return 1
	} else {
		return 0
	}
}

func (this *MnsModule) Watch() {
	client := ali_mns.NewAliMNSClient(this.Config.GetStringItem("Url"),
		this.Config.GetStringItem("AccessKeyId"),
		this.Config.GetStringItem("AccessKeySecret"))

	queue := ali_mns.NewMNSQueue(this.Config.GetStringItem("Queue"), client)

	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case resp := <-respChan:
				{
					//logs.Pretty("message:", string(resp.MessageBody))

					msg := &MNSmsg{}
					err := json.Unmarshal(resp.MessageBody, msg)
					if err != nil {
						fmt.Println(err)
						this.Warning(err.Error())
						break
					}

					content, err := base64.StdEncoding.DecodeString(msg.Payload)
					if err != nil {
						this.Warning(err.Error())
						break
					}

					fmt.Println("Topic: ", msg.Topic, "MessageType: ", msg.MessageType)
					fmt.Println("Message: ", msg.Messageid)
					fmt.Println("content: ", string(content))

					if this.Config.GetBoolItem("Delete") {
						if e := queue.DeleteMessage(resp.ReceiptHandle); e != nil {
							logs.Error(e)
						}
					}
				}

			case err := <-errChan:
				{
					if !strings.Contains(err.Error(), "code: MessageNotExist,") {
						logs.Error(err)
					}
				}
			}
		}
	}()

	queue.ReceiveMessage(respChan, errChan)

}

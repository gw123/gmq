package mnsModule

import (
	"github.com/gogap/ali_mns"
	"github.com/gogap/logs"
	"strings"
		"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
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
	client := ali_mns.NewAliMNSClient(this.Config.GetItem("Url"),
		this.Config.GetItem("AccessKeyId"),
		this.Config.GetItem("AccessKeySecret"))

	queue := ali_mns.NewMNSQueue(this.Config.GetItem("Queue"), client)

	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		for {
			select {
			case resp := <-respChan:
				{
					logs.Pretty("message:", string(resp.MessageBody))
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

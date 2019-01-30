package mqttModule

import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/common/common_types"
	"github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"time"
)

type MqttModule struct {
	base.BaseModule
	IotInstance     *Iot
	IsConnectAliIot bool
}

func NewMqttModule() *MqttModule {
	this := new(MqttModule)
	return this
}

//订阅 mqtt_upload  mqtt_log
func (this *MqttModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	//this.SetApp(app)
	this.BaseModule.Init(app, config)
	app.Sub("reply", this)
	app.Sub("sendMqttMsg", this)
	this.initAliIot(app, config)
	return nil
}

func (this *MqttModule) initAliIot(app interfaces.App, config interfaces.ModuleConfig) {
	params := Params{}
	params.ProductKey = config.GetItem("productKey")
	params.DeviceName = config.GetGlobalItem("deviceName")
	params.DeviceSecret = config.GetItem("deviceSecret")
	//this.Debug(fmt.Sprintf("PK:%s ,DN:%s ,DS:%s", params.ProductKey, params.DeviceName, params.DeviceSecret))
	params.OnConnectHandler = func(client mqtt.Client) {
		this.Info("连接阿里云成功")
		this.IotInstance.PublishInform(app.GetVersion())
		type VersionEvent struct {
			Data      string `json:"data"`
			EventName string `json:"event"`
		}

		versionEvent := VersionEvent{}
		versionEvent.Data = app.GetVersion()
		versionEvent.EventName = "version"

		data, err := json.Marshal(versionEvent)
		if err != nil {
			this.Error("json.Marshal:" + err.Error())
		} else {
			this.IotInstance.PublishRaw(data)
		}
		this.IsConnectAliIot = true
	}

	params.ConnectionLostHandler = func(client mqtt.Client, e error) {
		this.Info("和阿里云连接断开")
		this.IsConnectAliIot = false
	}

	params.DefaultHandel = func(client mqtt.Client, msg mqtt.Message) {
		switch msg.Topic() {
		case "/ota/device/upgrade/" + params.ProductKey + "/" + params.DeviceName:
			event := &common_types.Event{}
			json.Unmarshal(msg.Payload(), event)
			this.App.Pub(event)
			break;
		case "/" + params.ProductKey + "/" + params.DeviceName + "/get":
			event := &common_types.Event{}
			json.Unmarshal(msg.Payload(), event)
			this.App.Pub(event)
			break;
		default:
			this.Warning("Default TOPIC :" + msg.Topic() + " : " + string(msg.Payload()))
		}
	}
	params.Logger = this
	params.App = this.App
	this.IotInstance = NewIot(params)
	this.IotInstance.Connect()
	this.IotInstance.SubscribeGet()
	this.IotInstance.SubscribeUpgrade()
}

func (this *MqttModule) UnInit() error {
	this.App.UnSub("reply", this)
	this.App.UnSub("sendMqttMsg", this)
	return nil
}

func (this *MqttModule) Push(event interfaces.Event) (err error) {
	err = this.BaseModule.Push(event)
	return
}

func (this *MqttModule) GetStatus() uint64 {
	if this.IsConnectAliIot {
		return 1
	} else {
		return 0
	}
}

func (this *MqttModule) Start() {
	for ; ; {
		event := this.BaseModule.Pop()

		err := this.service(event)

		if err != nil {
			//执行失败
			//replay := common.NewResultEvent([]byte("失败" + err.Error()))
			//this.App.Pub(replay)
			this.Warning("service 执行失败 " + err.Error())
		} else {
			////执行成功
			//replay := common.NewResultEvent([]byte("成功"))
			//this.App.Pub(replay)
		}
		time.Sleep(time.Second)
	}
}

func (this *MqttModule) service(event interfaces.Event) error {
	this.Debug("service " + event.GetEventName() + " " + event.GetSourceModule() + " " + event.GetMsgId() + " " + string(event.GetPayload()))
	//data, err := json.Marshal()
	//if err != nil {
	//	this.Error("service " + err.Error())
	//	return nil
	//}

	this.IotInstance.PublishRaw(event.GetPayload())
	return nil
}

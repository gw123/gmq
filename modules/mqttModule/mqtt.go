package mqttModule

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"encoding/json"
	"os"
	"sync"
	erp_interfaces "github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/common/common_types"
	"github.com/fpay/erp-client-s/common"
)

const MaxMsgLen = 100

type MsgIds struct {
	Ids  []string
	head int
	tail int
}

func NewMsgIds() *MsgIds {
	this := new(MsgIds)
	this.head = 0
	this.Ids = make([]string, MaxMsgLen)
	return this
}

var msgIdsLock = sync.Mutex{}

func (msgIds *MsgIds) Check(msgId string) bool {
	msgIdsLock.Lock()
	defer msgIdsLock.Unlock()
	for _, id1 := range msgIds.Ids {
		if id1 == msgId {
			return false
		}
	}

	if msgIds.head >= MaxMsgLen {
		msgIds.head = 0
	}
	msgIds.Ids[msgIds.head] = msgId
	msgIds.head++;
	return true
}

var msgIds *MsgIds

func init() {
	msgIds = NewMsgIds()
}

type Device struct {
	DeviceName   string
	ProductKey   string
	DeviceSecret string
	IsLogin      bool
	Status       uint
	LasUsedtime  int64
}

type Iot struct {
	Host       string
	DeviceName string
	ProductKey string
	ClientId   string
	Username   string
	Password   string
	Sign       string
	Conn       mqtt.Client
	logOut     erp_interfaces.ModuleLogger
	App        erp_interfaces.App
	SubDevices []Device
}

type Params struct {
	ProductKey            string
	DeviceName            string
	DeviceSecret          string
	OnConnectHandler      mqtt.OnConnectHandler
	ConnectionLostHandler mqtt.ConnectionLostHandler
	Logger                erp_interfaces.ModuleLogger
	App                   erp_interfaces.App
	DefaultHandel         mqtt.MessageHandler
}

func NewIot(params Params) (iot *Iot) {
	iot = new(Iot)
	iot.SubDevices = make([]Device, 0)
	sign, timestamp := iot.GetSign(params.ProductKey, params.DeviceName, params.DeviceSecret)
	iot.Password = sign
	iot.ClientId = params.DeviceName + "|securemode=3,signmethod=hmacsha1,timestamp=" + timestamp + "|"
	iot.Username = params.DeviceName + "&" + params.ProductKey
	iot.DeviceName = params.DeviceName
	iot.ProductKey = params.ProductKey
	iot.Host = params.ProductKey + ".iot-as-mqtt.cn-shanghai.aliyuncs.com:1883"

	opts := mqtt.NewClientOptions().AddBroker(iot.Host).SetClientID(iot.ClientId).SetUsername(iot.Username).SetPassword(iot.Password)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(params.ConnectionLostHandler)
	opts.SetDefaultPublishHandler(params.DefaultHandel)
	opts.SetOnConnectHandler(params.OnConnectHandler)
	opts.SetMaxReconnectInterval(2 * time.Minute)
	iot.App = params.App
	iot.logOut = params.Logger
	iot.Conn = mqtt.NewClient(opts)
	return
}

func (this *Iot) Connect() (err error) {
	c := this.Conn
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (this *Iot) Close() {
	if this.Conn.IsConnected() {
		this.Conn.Disconnect(256)
	}
}

func (this *Iot) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	this.Conn.Publish(topic, qos, retained, payload)
	return nil
}

func (this *Iot) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := this.Conn.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		callback(client, message)
	}); token.Wait() && token.Error() != nil {
		this.writeLog("error", "Token publish: "+token.Error().Error())
		return token.Error()
	}
	return nil
}

func (this *Iot) SubscribeAndCheck(topic string, qos byte) error {
	if token := this.Conn.Subscribe(topic, qos, this.SubscribeGetCallback); token.Wait() && token.Error() != nil {
		this.writeLog("error", "Token publish: "+token.Error().Error())
		return token.Error()
	}
	return nil
}

func (this *Iot) SubscribeGetCallback(client mqtt.Client, message mqtt.Message) {
	topic := "/" + this.ProductKey + "/" + this.DeviceName + "/get"
	msg := &AliMsg{}
	err := json.Unmarshal(message.Payload(), msg)
	if err != nil {
		this.writeLog("error", "Topic "+topic+" 消息解密失败 "+err.Error()+" Payload: "+string(message.Payload()))
		return
	}
	if !msgIds.Check(msg.MsgId) {
		this.writeLog("warning", "msgId "+msg.MsgId+" Topic"+topic+" 重复消息")
		return
	}
	event := common_types.NewEvent(msg.MsgId ,message.Payload())
	this.App.Pub(event)
}

/***
 * 订阅get消息
 */
func (this *Iot) SubscribeGet() {
	topic := "/" + this.ProductKey + "/" + this.DeviceName + "/get"
	this.SubscribeAndCheck(topic, 1)
}

/***
 * 子设备
 */
func (this *Iot) SubscribeSubGet(subProductKey, subDeviceName string) {
	topic := "/" + subProductKey + "/" + subDeviceName + "/get"
	this.SubscribeAndCheck(topic, 0)
}

/***
 * 子设备注册
 */
func (this *Iot) PublishSubRegister(subProductKey, subDeviceName string) {
	data := "{'id': '%s', 'version':'1.0','params':[{'deviceName':'%s','productKey':'%s'}],'method':'thing.sub.register'}"
	data = fmt.Sprintf(data, this.getMsgId(), subDeviceName, subProductKey)
	topic := "/sys/" + this.ProductKey + "/" + this.DeviceName + "/thing/sub/register"
	this.Publish(topic, 0, false, []byte(data))
}

func (this *Iot) SubscribeSubRegisterReply() {
	topic_reply := "/sys/" + this.ProductKey + "/" + this.DeviceName + "/thing/sub/register_reply"
	this.Subscribe(topic_reply, 1, func(client mqtt.Client, message mqtt.Message) {
		msg, err := common.ParseAliMsg(message.Payload())
		if err != nil {
			this.writeLog("error", "SubRegister_reply json内容解析失败 "+string(message.Payload()))
			return
		}

		if msg.Code != 200 {
			this.writeLog("error", "SubRegister_reply 子设备注册失败 "+msg.Message)
			return
		}

		v, ok := msg.Data.([]interface{})
		if !ok {
			this.writeLog("error", "SubRegister_reply json内容解析失败->data解析失败 "+string(message.Payload()))
			return
		}
		for _, deviceData := range v {
			deviceInfo, ok := deviceData.(map[string]interface{})
			if !ok {
				this.writeLog("error", "SubRegister_reply json内容解析失败->data解析失败->不能转为map"+string(message.Payload()))
				continue
			}
			deviceSecret, _ := deviceInfo["deviceSecret"].(string)
			productKey, _ := deviceInfo["productKey"].(string)
			deviceName, _ := deviceInfo["deviceName"].(string)
			this.writeLog("info", "SubRegister_reply 注册成功: "+deviceName)
			go this.SubDeviceLogin(productKey, deviceName, deviceSecret)
		}
	})
}

func (this *Iot) SubDeviceLogin(productKey, deviceName, deviceSecret string) {
	this.AppendSubDevice(productKey, deviceName, deviceSecret)
	this.PublishSubAdd(productKey, deviceName, deviceSecret)
	time.Sleep(time.Second * 4)
	this.PublishSubLogin(productKey, deviceName, deviceSecret)
}

/***
 * 添加子设备
 */
func (this *Iot) PublishSubAdd(subProductKey, subDeviceName, subDeviceSecret string) {
	sign, timestamp := this.GetSign(subProductKey, subDeviceName, subDeviceSecret)
	data := `{"id":"%s","version":"1.0","params":[{"productKey" : "%s","deviceName" : "%s","clientId":"%s","sign":"%s","signmethod":"hmacSha1","timestamp":"%s"}],"method":"thing.topo.add"}`
	data = fmt.Sprintf(data, this.getMsgId(), subProductKey, subDeviceName, subDeviceName, sign, timestamp)
	topic := "/sys/" + this.ProductKey + "/" + this.DeviceName + "/thing/topo/add"
	this.SubscribeSubAddReply()
	this.Publish(topic, 0, true, []byte(data))
}

func (this *Iot) SubscribeSubAddReply() {
	topic_reply := "/sys/" + this.ProductKey + "/" + this.DeviceName + "/thing/topo/add_reply"
	this.Subscribe(topic_reply, 0, func(client mqtt.Client, message mqtt.Message) {
		msg, err := common.ParseAliMsg(message.Payload())
		if err != nil {
			this.writeLog("error", "PublishSubAdd "+"JSON解析失败")
			return
		}

		if msg.Code != 200 {
			this.writeLog("error", "PublishSubAdd "+msg.Message)
		}

		this.writeLog("info", "PublishSubAdd 子设备拓扑添加成功")
		return
	})
}

/***
 step [1, 100] -1：代表升级失败 -2：代表下载失败 -3：代表校验失败 -4：代表烧写失败
 desc 进度信息
 */
func (this *Iot) PublishProgress(step int8, desc string) {
	topic := "/ota/device/progress/" + this.ProductKey + "/" + this.DeviceName
	data := `{ "id": "%s", "params": {"step":"%d", "desc":" %s"}}`
	data = fmt.Sprintf(data, this.getMsgId(), step, desc)
	this.Publish(topic, 1, false, []byte(data))
}

func (this *Iot) SubscribeUpgrade() {
	topic := "/ota/device/upgrade/" + this.ProductKey + "/" + this.DeviceName
	this.Subscribe(topic, 1, this.SubscribeUpgradeCallback)
}

func (this *Iot) SubscribeUpgradeCallback(client mqtt.Client, message mqtt.Message) {
	//fmt.Println("SubscribeUpgradeCallback", message.Topic(), string(message.Payload()))
	update := common.UpdateResponse{}
	err := json.Unmarshal(message.Payload(), &update)
	if err != nil {
		this.writeLog("error", "SubscribeUpgrade"+"Json fail "+err.Error())
		return
	}
	if update.Message != "success" {
		this.writeLog("error", "SubscribeUpgrade "+update.Message)
		return
	}
	this.SyncUpgradeFile(update)
}

/***
 * 上报设备版本信息
 */
func (this *Iot) PublishInform(version string) {
	data := `{"id": "%s","params": {"version": "%s"}}`
	data = fmt.Sprintf(data, this.getMsgId(), version)
	topic := "/ota/device/inform/" + this.ProductKey + "/" + this.DeviceName
	this.Publish(topic, 0, false, []byte(data))
}

/***
 * 子设备登陆 这个函数会一直执行 ,所以在运行时要加上 go iot.PublishSubLogin
 */
func (this *Iot) PublishSubLogin(subProductKey, subDeviceName, subDeviceSecret string) {
	sign, timestamp := this.GetSign(subProductKey, subDeviceName, subDeviceSecret)
	data := `{"id":"%s","params":{"productKey":"%s","deviceName":"%s","clientId":"%s","sign":"%s","timestamp":"%s","signMethod":"hmacSha1","cleanSession":"false"}}`
	data = fmt.Sprintf(data, "ababab", subProductKey, subDeviceName, subDeviceName, sign, timestamp)
	topic := "/ext/session/" + this.ProductKey + "/" + this.DeviceName + "/combine/login"
	this.Publish(topic, 1, true, []byte(data))
}

/***
 * 子设备登陆回调函数
 */
func (this *Iot) SubscribeSubLoginReply() {
	topic_reply := "/ext/session/" + this.ProductKey + "/" + this.DeviceName + "/combine/login_reply"
	this.Subscribe(topic_reply, 1, func(client mqtt.Client, message mqtt.Message) {
		msg := common.LoginResponse{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			this.writeLog("error", "SubLogin_reply Json 解析失败"+string(message.Payload()))
			return
		}

		if msg.Code != 200 {
			this.writeLog("error", "SubLogin_reply Json 登陆失败"+msg.Message)
			return
		}
		this.writeLog("info", "SubLogin_reply "+msg.Data.DeviceName+" 登陆成功"+msg.Message)

		/*订阅主题*/
		this.SubscribeSubGet(msg.Data.ProductKey, msg.Data.DeviceName)
	})
}

/***
 * 子设备下线
 */
func (this *Iot) PublishSubLoginOut(subProductKey, subDeviceName string) {
	data := `{"id":"%s","params":{"productKey":"%s","deviceName":"%s",}}`
	data = fmt.Sprintf(data, this.getMsgId(), subProductKey, subDeviceName)
	topci := "/" + this.ProductKey + "/" + this.DeviceName + "/combine/logout"
	topci_reply := "/" + this.ProductKey + "/" + this.DeviceName + "/combine/logout_reply"
	this.Publish(topci, 1, false, []byte(data))
	this.Subscribe(topci_reply, 1, func(client mqtt.Client, message mqtt.Message) {
		msg, err := common.ParseAliMsg(message.Payload())
		if err != nil {
			this.writeLog("error", "SubLoginOut :"+subDeviceName+" "+err.Error())
			return
		}
		if msg.Code != 200 {
			this.writeLog("error", "SubLoginOut :"+subDeviceName+" "+msg.Message)
			return
		}

	})
}

/***
 * 计算签名
 */
func (this *Iot) GetSign(productKey, deviceName, deviceSecret string) (string, string) {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	str := "clientId" + deviceName + "deviceName" + deviceName + "productKey" + productKey + "timestamp" + timestamp;
	key := []byte(deviceSecret)
	hmacHandel := hmac.New(sha1.New, key)
	hmacHandel.Write([]byte(str))
	res := hmacHandel.Sum(nil)
	return hex.EncodeToString(res), timestamp
}

/***
 * 获取一个唯一的消息Id
 */
func (this *Iot) getMsgId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func (this *Iot) SetLogOutPut(writer erp_interfaces.ModuleLogger) {
	this.logOut = writer
}

func (this *Iot) writeLog(logType, Content string) {
	switch logType {
	case "warning":
		this.logOut.Warning(Content)
		break;
	case "info":
		this.logOut.Info(Content)
		break
	case "error":
		this.logOut.Error(Content)
		break;
	case "debug":
		this.logOut.Debug(Content)
		break
	default:
		this.logOut.Info(Content)
		break
	}
}

func (this *Iot) Write(data []byte) (int, error) {
	if this.Conn.IsConnected() {
		this.PublishRaw(data)
	}
	return 0, nil
}

var appendSubDevicesMutex = sync.Mutex{}

/***
 * 添加子设备
 */
func (this *Iot) AppendSubDevice(subProductKey, subDeviceName, subDeviceSecret string) (Device) {
	subDevice := Device{}
	subDevice.ProductKey = subProductKey
	subDevice.DeviceName = subDeviceName
	subDevice.DeviceSecret = subDeviceSecret
	appendSubDevicesMutex.Lock()
	this.SubDevices = append(this.SubDevices, subDevice)
	appendSubDevicesMutex.Unlock()

	return subDevice
}

func (this *Iot) PublishRaw(data []byte) {
	topic := "/" + this.ProductKey + "/" + this.DeviceName + "/update"
	this.Publish(topic, 1, false, data)
}

func (this *Iot) PublishLog(log []byte) {
	topic := "/" + this.ProductKey + "/" + this.DeviceName + "/update"

	type Log struct {
		Timestamp int64  `json:"timestamp"`
		Event     string `json:"event"`
		Data      string `json:"data"`
	}

	logData := Log{}
	logData.Timestamp = time.Now().Unix()
	logData.Data = string(log)
	logData.Event = "log"
	data, err := json.Marshal(logData)
	if err != nil {
		return
	}
	this.Publish(topic, 1, false, data)
}

func (this *Iot) SyncUpgradeFile(data common.UpdateResponse) {
	fileContent := `%s
%s
%s`
	fileContent = fmt.Sprintf(fileContent, data.Data.Md5, data.Data.Url, data.Data.Version)
	var file *os.File;
	var err error;
	flag := false
	for i := 5; i > 0; i-- {
		file, err = os.OpenFile("upgrade.plan", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
		if err != nil {
			continue
		}
		flag = true
		time.Sleep(time.Second)
		break
	}

	if !flag {
		this.writeLog("error", "SyncUpgradeFile： "+err.Error())
		return
	}

	defer file.Close()
	file.Write([]byte(fileContent))
}

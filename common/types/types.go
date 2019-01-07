package types

import (
	"encoding/json"
	"time"
	"strconv"
)

type AliMsg struct {
	Id      string      `json:"id"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type LhMsg struct {
	Timestamp   int64    `json:"timestamp"`
	Expired     int      `json:"expired"`
	MsgId       string   `json:"msgId"`
	DeviceNames []string `json:"deviceNames"`
	EventName       string   `json:"event"`
	Payload        string   `json:"data"`
}

func (this *LhMsg) GetMsgId() string {
	return this.MsgId
}

func (this *LhMsg) GetEventName() string {
	return this.EventName
}

func (this *LhMsg) GetPayload() []byte {
	return []byte(this.Payload)
}


type Register struct {
	IotId        string `json:"iotId"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

type LoginResponse struct {
	Id      string `json:"id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data struct {
		ProductKey string `json:"productKey"`
		DeviceName string `json:"deviceName"`
	} `json:"data"`
}

type ChangeIpResponse struct {
	Timestamp int64  `json:"timestamp"`
	MsgId     string `json:"msgId"`
	Event     string `json:"event"`
	Data struct {
		Status int64  `json:"status"`
		Ip     string `json:"ip"`
		Msg    string `json:"msg"`
	} `json:"data"`
}

func NewChangeIpResponse(status int64, ip string) *ChangeIpResponse {
	this := new(ChangeIpResponse)
	this.Timestamp = time.Now().Unix()
	this.Event = "changeIp"
	this.MsgId = strconv.Itoa(int(time.Now().UnixNano()))
	this.Data.Status = status
	this.Data.Ip = ip
	return this
}

type UpdateResponse struct {
	Code int `json:"code"`
	Data struct {
		Size    int    `json:"size"`
		Url     string `json:"url"`
		Md5     string `json:"md5"`
		Version string `json:"version"`
	} `json:"data"`
	Message string `json:"message"`
}

/***
 *
 */
func ParseAliMsg(payload []byte) (msg AliMsg, err error) {
	err = json.Unmarshal(payload, &msg)
	if err != nil {
		return msg, err
	}
	return msg, err
}

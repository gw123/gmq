package mqttModule

import (
	"encoding/json"
)

type AliMsg struct {
	MsgId   string      `json:"id"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}


type Register struct {
	IotId        string `json:"iotId"`
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
}

type LoginResponse struct {
	MsgId      string `json:"id"`
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

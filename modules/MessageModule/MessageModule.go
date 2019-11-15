package MessageModule

import (
	"errors"
	"fmt"
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
)

type MessageModule struct {
	base.BaseModule
	client *AliyunClient
}

func NewMessageModule() *MessageModule {
	this := new(MessageModule)

	return this
}

func (this *MessageModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	providers := this.BaseModule.GetConfig().GetMapItem("providers")
	deft, ok := providers["default"].(map[string]interface{})
	if !ok {
		return errors.New("providers 配置格式有问题,缺少default配置")
	}
	keyId := deft["key_id"].(string)
	keySecret := deft["key_secret"].(string)
	regionId := deft["region_id"].(string)
	signName := deft["sign_name"].(string)
	templateCode := deft["template_code"].(string)
	client, err := NewAlyunClient(regionId, keyId, keySecret, signName, templateCode)
	if err != nil {
		return err
	}
	this.client = client
	//在init 之后才能订阅消息
	this.Sub("sendMobileMessage")
	return nil
}

func (this *MessageModule) Handle(event interfaces.Msg) error {
	mobileMessageEvent := event.(*gmsg.MobileMessageEvent)
	content := fmt.Sprintf(`{"code":"%s"}`, mobileMessageEvent.Code)
	return this.client.SendCode(mobileMessageEvent.Modbile, content)
}

func (this *MessageModule) Watch(index int) {
	//if index%60 == 0 {
	//	this.Debug("发送短信")
	//	e := common_types.NewMobileMessageEvent("1234", "mo")
	//	this.Push(e)
	//}
	return
}

package MessageModule

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/pkg/errors"
)

type AliyunClient struct {
	client       *dysmsapi.Client
	module       interfaces.Module
	SignName     string
	TemplateCode string
}

func NewAlyunClient(region, key, secret, signName, templateCode string) (*AliyunClient, error) {
	client, err := dysmsapi.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		return nil, err
	}

	return &AliyunClient{
		SignName:     signName,
		TemplateCode: templateCode,
		client:       client,
	}, nil
}

func (c *AliyunClient) SendCode(mobile, content string) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = c.SignName
	request.TemplateCode = c.TemplateCode
	request.TemplateParam = content

	response, err := c.client.SendSms(request)
	if err != nil {
		return err
	}
	if response.Message != "OK" {
		return errors.New(response.RequestId + "," + response.Message)
	}
	return nil
}

package serverNodeModule

import (
	"github.com/gw123/ghelper/map-tool"
	"net/url"
	"strings"
	"fmt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

/*
http://iot.cn-shanghai.aliyuncs.com/?
MessageContent=aGVsbG93b3JsZA%3D&
Action=Pub&
Timestamp=2017-10-02T09%3A39%3A41Z&
SignatureVersion=1.0&
ServiceCode=iot&
Format=XML&
Qos=0&
SignatureNonce=0715a395-aedf-4a41-bab7-746b43d38d88&
Version=2017-04-20&
AccessKeyId=testid&
SignatureMethod=HMAC-SHA1&
RegionId=cn-shanghai&
ProductKey=12345abcdeZ&
TopicFullName=%2FproductKey%2Ftestdevice%2Fget
*/

type CommonParams struct {
	Format      string
	Version     string
	AccessKeyId string
	//HMAC-SHA1
	SignatureMethod  string
	Timestamp        string
	SignatureVersion string
	SignatureNonce   string
	RegionId         string
}

type SignHttpRequest struct {
	Host     string
	Protocol string
	Ip       string
	Path     string

	commonParams    CommonParams
	AccessKeySecret string
	BaseUrl         string
}

func NewClient(AccessKeyId, AccessKeySecret string) *SignHttpRequest {
	params := CommonParams{
		Format:           "JSON",
		Version:          "2017-04-20",
		AccessKeyId:      AccessKeyId,
		SignatureMethod:  "HMAC-SHA1",
		SignatureVersion: "1.0",
		RegionId:         "cn-shanghai",
	}

	client := &SignHttpRequest{
		commonParams:    params,
		AccessKeySecret: AccessKeySecret,
	}

	return client
}

func (this *SignHttpRequest) SetFormat(format string) {
	this.commonParams.Format = format
}

func (this *SignHttpRequest) SetRegionId(regionId string) {
	this.commonParams.RegionId = regionId
}

func (this *SignHttpRequest) Post(url string, params interface{}, response interface{}) error {
	method := "POST"
	sign, queryStr := this.makeSign(method, params)
	//url := fmt.Sprintf("%s://%s", this.Protocol, this.Host)
	queryStr = queryStr + "&Signature=" + sign
	// Get
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Post(url, "application/x-www-form-urlencoded", strings.NewReader(queryStr))

	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, response)
	if err != nil {
		return err
	}
	return nil
}

/***
 * 计算签名
 */
func (this *SignHttpRequest) makeSign(method string, params interface{}) (string, string) {
	this.commonParams.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	this.commonParams.SignatureNonce = fmt.Sprintf("%d", time.Now().Unix())

	requestParams := map_tool.MergeField(this.commonParams, params)
	sortKeys := map_tool.SortMapKeys(requestParams)
	mixVals := make([]string, 0)

	for _, key := range sortKeys {
		val := requestParams[key]
		mixVal := url.QueryEscape(key) + "=" + url.QueryEscape(val)
		mixVals = append(mixVals, mixVal)
	}
	mixValsStr := strings.Join(mixVals, "&")
	toSignStr := method + "&%2F&" + url.QueryEscape(mixValsStr)

	hmacHandel := hmac.New(sha1.New, []byte(this.AccessKeySecret+"&"))
	hmacHandel.Write([]byte(toSignStr))
	res := hmacHandel.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(res)

	return sign, mixValsStr
}

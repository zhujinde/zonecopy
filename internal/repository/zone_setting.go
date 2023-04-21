package repository

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"zonecopy/internal/domain/entity"
)

const (
	StatusOn  = "on"
	StatusOff = "off"
)

// ZoneSetting 站点配置管理。
type ZoneSetting struct {
	Account *entity.AccountBaseInfo
	Request *teo.ModifyZoneSettingRequest
}

func NewZoneSetting(a *entity.AccountBaseInfo) *ZoneSetting {
	return &ZoneSetting{
		Account: a,
		Request: teo.NewModifyZoneSettingRequest(),
	}
}

func (z *ZoneSetting) OpenClientIpCountry(zoneId, head string) error {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := teo.NewModifyZoneSettingRequest()
	request.ZoneId = common.StringPtr(zoneId)
	request.ClientIpCountry = &teo.ClientIpCountry{
		Switch:     common.StringPtr("on"),
		HeaderName: common.StringPtr(head),
	}
	body, _ := json.Marshal(request)
	log.Printf("###### API: OpenClientIpCountry Request: %#v", string(body))

	// 返回的resp是一个ModifyZoneSettingResponse的实例，与请求对象对应
	response, err := client.ModifyZoneSetting(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: OpenClientIpCountry response: %#v\n", response.ToJsonString())
	return nil
}

func (z *ZoneSetting) ModifyHttp2(zoneId string, status bool) error {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := teo.NewModifyZoneSettingRequest()
	request.ZoneId = common.StringPtr(zoneId)
	request.Https = &teo.Https{
		Http2: common.StringPtr(StatusOff),
	}
	if status {
		request.Https.Http2 = common.StringPtr(StatusOn)
	}
	body, _ := json.Marshal(request)
	log.Printf("###### API: ModifyHttp2 Request: %#v", string(body))

	// 返回的resp是一个ModifyZoneSettingResponse的实例，与请求对象对应
	response, err := client.ModifyZoneSetting(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return err
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: ModifyHttp2 response: %#v\n", response.ToJsonString())
	return nil
}

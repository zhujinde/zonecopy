package repository

import (
	"fmt"
	"log"

	"github.com/mulinbc/zerr"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"zonecopy/internal/domain/entity"
)

// ZoneSettingManager 站点加速配置。
type ZoneSettingManager struct {
	Account *entity.AccountBaseInfo
}

func NewZoneSettingManager(a *entity.AccountBaseInfo) *ZoneSettingManager {
	return &ZoneSettingManager{
		Account: a,
	}
}

func (z *ZoneSettingManager) DescribeZoneSetting(zoneId string) (*teo.ZoneSetting, error) {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)
	request := teo.NewDescribeZoneSettingRequest()
	request.ZoneId = common.StringPtr(zoneId)
	response, err := client.DescribeZoneSetting(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "internal error")
	}
	log.Printf("[API] DescribeZoneSetting response: %#v", response.ToJsonString())
	return response.Response.ZoneSetting, nil
}

func (z *ZoneSettingManager) ModifyZoneSetting(request *teo.ModifyZoneSettingRequest) error {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)
	log.Printf("[API] ModifyZoneSetting Request: %#v", request.ToJsonString())
	response, err := client.ModifyZoneSetting(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return zerr.Wrap(err, "internal error")
	}
	log.Printf("[API] ModifyZoneSetting response: %#v", response.ToJsonString())
	return nil
}

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

// OriginManager 源站导入。
type OriginManager struct {
	Account *entity.AccountBaseInfo
}

func NewOriginManager(a *entity.AccountBaseInfo) *OriginManager {
	return &OriginManager{
		Account: a,
	}
}

func (o *OriginManager) DescribeOriginGroupList(zoneId string) ([]*teo.OriginGroup, error) {
	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)

	request := teo.NewDescribeOriginGroupRequest()
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.Filters = []*teo.AdvancedFilter{
		&teo.AdvancedFilter{
			Name:   common.StringPtr("zone-id"),
			Values: common.StringPtrs([]string{zoneId}),
		},
	}
	log.Printf("[API] DescribeOriginGroupList Request: %#v", request.ToJsonString())

	response, err := client.DescribeOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "interal error")
	}
	log.Printf("[API] DescribeOriginGroupList response: %#v", response.ToJsonString())

	return response.Response.OriginGroups, nil
}

func (o *OriginManager) GetOriginIdByName(zoneId, groupName string) (string, error) {
	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)

	request := teo.NewDescribeOriginGroupRequest()
	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.Filters = []*teo.AdvancedFilter{
		&teo.AdvancedFilter{
			Name:   common.StringPtr("zone-id"),
			Values: common.StringPtrs([]string{zoneId}),
		},
		&teo.AdvancedFilter{
			Name:   common.StringPtr("origin-group-name"),
			Values: common.StringPtrs([]string{groupName}),
		},
	}
	log.Printf("[API] IsOriginExist Request: %#v", request.ToJsonString())

	response, err := client.DescribeOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return "", zerr.Wrap(err, "interal error")
	}
	log.Printf("[API] IsOriginExist response: %#v", response.ToJsonString())
	if *response.Response.TotalCount == 0 {
		return "", nil
	}
	if *response.Response.TotalCount != 1 {
		return "", fmt.Errorf("abnormal response")
	}
	return *response.Response.OriginGroups[0].OriginGroupId, nil
}

func (o *OriginManager) CreateOrigin(request *teo.CreateOriginGroupRequest) (string, error) {
	id, err := o.GetOriginIdByName(*request.ZoneId, *request.OriginGroupName)
	if err != nil {
		return "", zerr.Wrap(err, "o.IsOriginExist failed")
	}
	if id != "" {
		return id, nil
	}

	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)
	log.Printf("[API] CreateOrigin Request: %#v", request.ToJsonString())

	response, err := client.CreateOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return "", zerr.Wrap(err, "interal error")
	}
	log.Printf("[API] CreateOrigin response: %#v", response.ToJsonString())
	return *response.Response.OriginGroupId, nil
}

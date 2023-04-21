package repository

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mulinbc/zerr"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"zonecopy/internal/domain/entity"
)

// OriginManager 源站导入管理。
type OriginManager struct {
	Account *entity.AccountBaseInfo
}

func NewOriginManager(a *entity.AccountBaseInfo) *OriginManager {
	return &OriginManager{
		Account: a,
	}
}

func (o *OriginManager) DescribeOriginGroupList(zoneId string) ([]*teo.OriginGroup, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := teo.NewDescribeOriginGroupRequest()

	request.Offset = common.Uint64Ptr(0)
	request.Limit = common.Uint64Ptr(10)
	request.Filters = []*teo.AdvancedFilter{
		&teo.AdvancedFilter{
			Name:   common.StringPtr("zone-id"),
			Values: common.StringPtrs([]string{zoneId}),
		},
	}
	body, _ := json.Marshal(request)
	log.Printf("###### API: DescribeOriginGroupList Request: %#v", string(body))

	// 返回的resp是一个DescribeOriginGroupResponse的实例，与请求对象对应
	response, err := client.DescribeOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "interal error")
	}
	log.Printf("###### API: DescribeOriginGroupList response: %#v\n", response.ToJsonString())

	return response.Response.OriginGroups, nil
}

func (o *OriginManager) GetOriginIdByName(zoneId, groupName string) (string, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
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
	body, _ := json.Marshal(request)
	log.Printf("###### API: IsOriginExist Request: %#v", string(body))

	// 返回的resp是一个DescribeOriginGroupResponse的实例，与请求对象对应
	response, err := client.DescribeOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return "", zerr.Wrap(err, "interal error")
	}
	log.Printf("###### API: IsOriginExist response: %#v\n", response.ToJsonString())

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

	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		o.Account.SecretId,
		o.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = o.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, o.Account.Region, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	body, _ := json.Marshal(request)
	log.Printf("###### API: CreateOrigin Request: %#v", string(body))

	// 返回的resp是一个CreateOriginGroupResponse的实例，与请求对象对应
	response, err := client.CreateOriginGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return "", zerr.Wrap(err, "interal error")
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: CreateOrigin response: %#v\n", response.ToJsonString())
	return *response.Response.OriginGroupId, nil
}

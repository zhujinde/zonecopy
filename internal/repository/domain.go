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

// DomainManager 域名导入管理。
type DomainManager struct {
	Account *entity.AccountBaseInfo
}

func NewDomainManager(a *entity.AccountBaseInfo) *DomainManager {
	return &DomainManager{
		Account: a,
	}
}

func (z *DomainManager) IsDomainExist(zoneId, host string) (bool, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
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
	request := teo.NewDescribeAccelerationDomainsRequest()

	request.ZoneId = common.StringPtr(zoneId)
	request.Filters = []*teo.AdvancedFilter{
		&teo.AdvancedFilter{
			Name:   common.StringPtr("domain-name"),
			Values: common.StringPtrs([]string{host}),
		},
	}
	body, _ := json.Marshal(request)
	log.Printf("###### API: IsDomainExist Request: %#v", string(body))

	// 返回的resp是一个DescribeAccelerationDomainsResponse的实例，与请求对象对应
	response, err := client.DescribeAccelerationDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return false, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return false, zerr.Wrap(err, "internal error")
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: IsDomainExist response: %#v\n", response.ToJsonString())
	if *response.Response.TotalCount > 0 {
		return true, nil
	}
	return false, nil
}

func (z *DomainManager) DescribeDomainListDetail(zoneId string) ([]*teo.AccelerationDomain, error) {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
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
	request := teo.NewDescribeAccelerationDomainsRequest()
	request.ZoneId = common.StringPtr(zoneId)
	body, _ := json.Marshal(request)
	log.Printf("###### API: DescribeDomainListDetail Request: %#v", string(body))

	// 返回的resp是一个DescribeAccelerationDomainsResponse的实例，与请求对象对应
	response, err := client.DescribeAccelerationDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "internal error")
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: DescribeDomainListDetail response: %#v\n", response.ToJsonString())
	return response.Response.AccelerationDomains, nil
}

func (z *DomainManager) CreateDomain(request *teo.CreateAccelerationDomainRequest) error {
	ok, err := z.IsDomainExist(*request.ZoneId, *request.DomainName)
	if err != nil {
		return zerr.Wrap(err, "z.IsDomainExist failed")
	}
	if ok {
		return nil
	}

	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
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
	//request := teo.NewCreateAccelerationDomainRequest()
	//request.ZoneId = common.StringPtr(zoneId)
	//request.DomainName = common.StringPtr(host)
	//request.OriginInfo = &teo.OriginInfo{
	//	OriginType: common.StringPtr("IP_DOMAIN"),
	//	Origin:     common.StringPtr(origin),
	//}

	body, _ := json.Marshal(request)
	log.Printf("###### API: CreateDomain Request: %#v", string(body))
	// 返回的resp是一个CreateAccelerationDomainResponse的实例，与请求对象对应
	response, err := client.CreateAccelerationDomain(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return zerr.Wrap(err, "internal error")
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: CreateDomain response: %#v\n", response.ToJsonString())
	return nil
}

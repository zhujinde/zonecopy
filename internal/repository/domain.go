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

// DomainManager 域名导入。
type DomainManager struct {
	Account *entity.AccountBaseInfo
}

func NewDomainManager(a *entity.AccountBaseInfo) *DomainManager {
	return &DomainManager{
		Account: a,
	}
}

func (z *DomainManager) IsDomainExist(zoneId, host string) (bool, error) {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)

	request := teo.NewDescribeAccelerationDomainsRequest()

	request.ZoneId = common.StringPtr(zoneId)
	request.Filters = []*teo.AdvancedFilter{
		&teo.AdvancedFilter{
			Name:   common.StringPtr("domain-name"),
			Values: common.StringPtrs([]string{host}),
		},
	}
	log.Printf("[API] IsDomainExist Request: %v", request.ToJsonString())

	response, err := client.DescribeAccelerationDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return false, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return false, zerr.Wrap(err, "internal error")
	}
	log.Printf("[API] IsDomainExist response: %#v\n", response.ToJsonString())
	if *response.Response.TotalCount > 0 {
		return true, nil
	}
	return false, nil
}

func (z *DomainManager) DescribeDomainListDetail(zoneId string) ([]*teo.AccelerationDomain, error) {
	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)

	request := teo.NewDescribeAccelerationDomainsRequest()
	request.ZoneId = common.StringPtr(zoneId)
	body, _ := json.Marshal(request)
	log.Printf("[API] DescribeDomainListDetail Request: %#v", string(body))

	response, err := client.DescribeAccelerationDomains(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "internal error")
	}

	log.Printf("[API] DescribeDomainListDetail response: %#v\n", response.ToJsonString())
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

	credential := common.NewCredential(
		z.Account.SecretId,
		z.Account.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = z.Account.EndPoint
	client, _ := teo.NewClient(credential, z.Account.Region, cpf)
	log.Printf("[API] CreateDomain Request: %v", request.ToJsonString())

	response, err := client.CreateAccelerationDomain(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return zerr.Wrap(err, "internal error")
	}

	log.Printf("[API] CreateDomain response: %#v", response.ToJsonString())
	return nil
}

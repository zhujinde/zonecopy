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

// RuleEngineManager 规则引擎管理。
type RuleEngineManager struct {
	Account *entity.AccountBaseInfo
}

func NewRuleEngineManager(a *entity.AccountBaseInfo) *RuleEngineManager {
	return &RuleEngineManager{
		Account: a,
	}
}

func (r *RuleEngineManager) DescribeRuleList(zoneId string) ([]*teo.RuleItem, error) {
	credential := common.NewCredential(
		r.Account.SecretId,
		r.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = r.Account.EndPoint
	client, _ := teo.NewClient(credential, r.Account.Region, cpf)
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := teo.NewDescribeRulesRequest()
	request.ZoneId = common.StringPtr(zoneId)

	body, _ := json.Marshal(request)
	log.Printf("###### API: DescribeOriginGroupList Request: %#v", string(body))

	// 返回的resp是一个DescribeRulesResponse的实例，与请求对象对应
	response, err := client.DescribeRules(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "interal error")
	}
	log.Printf("###### API: DescribeOriginGroupList response: %#v\n", response.ToJsonString())

	return response.Response.RuleItems, nil
}

func (r *RuleEngineManager) CreateRule(request *teo.CreateRuleRequest) error {
	rules, err := r.DescribeRuleList(*request.ZoneId)
	if err != nil {
		return fmt.Errorf("DescribeRuleList failed, err: %v\n", err)
	}
	for _, v := range rules {
		// 同名规则已存在，直接跳过
		if *request.RuleName == *v.RuleName {
			return nil
		}
	}

	body, _ := json.Marshal(request)
	fmt.Printf("###### API: CreateRule Request: %#v", string(body))

	credential := common.NewCredential(
		r.Account.SecretId,
		r.Account.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = r.Account.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := teo.NewClient(credential, r.Account.Region, cpf)

	// 返回的resp是一个CreateRuleResponse的实例，与请求对象对应
	response, err := client.CreateRule(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return fmt.Errorf("an abnormal error has returned: %v", err)
	}
	// 输出json格式的字符串回包
	log.Printf("###### API: CreateRule response: %#v\n", response.ToJsonString())
	return nil
}

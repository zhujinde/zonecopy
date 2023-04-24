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

// RuleEngineManager 规则引擎。
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
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = r.Account.EndPoint
	client, _ := teo.NewClient(credential, r.Account.Region, cpf)
	request := teo.NewDescribeRulesRequest()
	request.ZoneId = common.StringPtr(zoneId)
	log.Printf("[API] DescribeOriginGroupList Request: %#v", request.ToJsonString())

	response, err := client.DescribeRules(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an API error has returned: %s", err)
	}
	if err != nil {
		return nil, zerr.Wrap(err, "interal error")
	}
	log.Printf("[API] DescribeOriginGroupList response: %#v\n", response.ToJsonString())

	return response.Response.RuleItems, nil
}

func (r *RuleEngineManager) IsRuleExist(zoneId, ruleName string) (bool, error) {
	rules, err := r.DescribeRuleList(zoneId)
	if err != nil {
		return false, fmt.Errorf("DescribeRuleList failed, err: %v\n", err)
	}
	for _, v := range rules {
		// 同名规则已存在，直接跳过
		if ruleName == *v.RuleName {
			return true, nil
		}
	}

	return false, nil
}

func (r *RuleEngineManager) CreateRule(request *teo.CreateRuleRequest) error {
	if ok, err := r.IsRuleExist(*request.ZoneId, *request.RuleName); err != nil {
		return err
	} else if ok {
		log.Printf("rule name: %v is already exist \n", *request.RuleName)
		return nil
	}
	log.Printf("[API] CreateRule Request: %#v", request.ToJsonString())
	credential := common.NewCredential(
		r.Account.SecretId,
		r.Account.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = r.Account.EndPoint
	client, _ := teo.NewClient(credential, r.Account.Region, cpf)

	response, err := client.CreateRule(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("an API error has returned: %v", err)
	}
	if err != nil {
		return zerr.Wrap(err, "interal error")
	}
	log.Printf("[API] CreateRule response: %#v\n", response.ToJsonString())
	return nil
}

package usecase

import (
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"strings"
	"zonecopy/internal/domain/entity"
	"zonecopy/internal/repository"
)

// ZoneCopyManager 站点配置拷贝。
type ZoneCopyManager struct {
	config         *entity.ZoneCopyConfig
	originImporter *repository.OriginManager
	domainImporter *repository.DomainManager
	ruleImporter   *repository.RuleEngineManager

	isOriginInit   bool              // 标识以下两个源站组配置信息是否初始化了
	templateOrigin map[string]string // 旧的groupId -> groupName
	targetOrigin   map[string]string // 新的groupName -> groupId
}

func NewZoneCopyManager(c *entity.ZoneCopyConfig) *ZoneCopyManager {
	if c == nil || c.Account == nil {
		log.Println("empty config")
		return nil
	}
	return &ZoneCopyManager{
		config:         c,
		originImporter: repository.NewOriginManager(c.Account),
		domainImporter: repository.NewDomainManager(c.Account),
		ruleImporter:   repository.NewRuleEngineManager(c.Account),

		isOriginInit:   false,
		templateOrigin: make(map[string]string),
		targetOrigin:   make(map[string]string),
	}
}

// ImportOrigin 源站导入。
func (z *ZoneCopyManager) ImportOrigin() error {
	oldGroups, err := z.originImporter.DescribeOriginGroupList(z.config.TemplateZoneId)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range oldGroups {
		req := teo.NewCreateOriginGroupRequest()
		req.ZoneId = common.StringPtr(z.config.TargetZoneId)
		req.OriginType = v.OriginType
		req.OriginGroupName = v.OriginGroupName
		req.ConfigurationType = v.ConfigurationType
		req.OriginRecords = v.OriginRecords
		req.HostHeader = v.HostHeader
		_, err = z.originImporter.CreateOrigin(req)
		if err != nil {
			fmt.Printf("源站：%v 导入失败, err: %v\n", req.OriginGroupName, err)
			return err
		}
	}
	return nil
}

// ImportDomains 域名导入。
func (z *ZoneCopyManager) ImportDomains() error {
	oldDomains, err := z.domainImporter.DescribeDomainListDetail(z.config.TemplateZoneId)
	if err != nil {
		fmt.Println(err)
	}
	// TODO：验证对象存储源站是否正常
	for _, v := range oldDomains {
		req := teo.NewCreateAccelerationDomainRequest()
		req.ZoneId = common.StringPtr(z.config.TargetZoneId)
		req.DomainName = common.StringPtr(z.getNewDomainName(*v.DomainName))
		req.OriginInfo, err = z.converDomainOrigin(v.OriginDetail)
		if err != nil {
			fmt.Printf("domain：%v -> %v convert config failed， err: %v\n", *v.DomainName, *req.DomainName, err)
			return err
		}
		if err = z.domainImporter.CreateDomain(req); err != nil {
			fmt.Printf("domain：%v -> %v import failed， err: %v\n", *v.DomainName, *req.DomainName, err)
			return err
		}
	}
	return nil
}

// converDomainOrigin 域名导入时调整源站信息。
func (z *ZoneCopyManager) converDomainOrigin(old *teo.OriginDetail) (*teo.OriginInfo, error) {
	nw := &teo.OriginInfo{}
	nw.OriginType = old.OriginType
	nw.Origin = old.Origin
	nw.BackupOrigin = old.BackupOrigin
	nw.PrivateAccess = old.PrivateAccess
	nw.PrivateParameters = old.PrivateParameters

	// 源站组的话需替换OriginGroupId
	if *nw.OriginType != "ORIGIN_GROUP" {
		return nw, nil
	}
	if nw.Origin != nil && *nw.Origin != "" {
		id, err := z.getNewGroupId(*nw.Origin)
		if err != nil {
			return nil, err
		}
		nw.Origin = common.StringPtr(id)
	}
	if nw.BackupOrigin != nil && *nw.BackupOrigin != "" {
		id, err := z.getNewGroupId(*nw.BackupOrigin)
		if err != nil {
			return nil, err
		}
		nw.BackupOrigin = common.StringPtr(id)
	}
	return nw, nil
}

// getNewGroupId 旧站点GroupId转换为新站点GroupId。
func (z *ZoneCopyManager) getNewGroupId(old string) (string, error) {
	if !z.isOriginInit {
		oldGroups, err := z.originImporter.DescribeOriginGroupList(z.config.TemplateZoneId)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range oldGroups {
			z.templateOrigin[*v.OriginGroupId] = *v.OriginGroupName
		}
		newGroups, err := z.originImporter.DescribeOriginGroupList(z.config.TargetZoneId)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range newGroups {
			z.targetOrigin[*v.OriginGroupName] = *v.OriginGroupId
		}
		z.isOriginInit = true
	}
	name, ok := z.templateOrigin[old]
	if !ok {
		return "", fmt.Errorf("not find old origin name")
	}
	id, ok := z.targetOrigin[name]
	if !ok {
		return "", fmt.Errorf("not find new origin id")
	}
	return id, nil
}

// getNewDomainName 旧站点域名转换为新站点域名。
func (z *ZoneCopyManager) getNewDomainName(old string) string {
	// TODO: 域名转换方式会有小概率bug
	return strings.Replace(old, z.config.TemplateZone, z.config.TargetZone, 1)
}

// ImportRuleEngineRules 规则引擎中规则的导入。
func (z *ZoneCopyManager) ImportRuleEngineRules() error {
	oldRules, err := z.ruleImporter.DescribeRuleList(z.config.TemplateZoneId)
	if err != nil {
		fmt.Println(err)
	}
	// 逆序导入，最终展示保持和原站点一致
	l := len(oldRules)
	for i := l - 1; i >= 0; i-- {
		v := oldRules[i]
		req := teo.NewCreateRuleRequest()
		req.ZoneId = common.StringPtr(z.config.TargetZoneId)
		req.RuleName = v.RuleName
		req.Status = common.StringPtr("enable")
		req.Rules, err = z.convertRules(v.Rules)
		if err != nil {
			fmt.Printf("rule name: %v convert config failed, err: %v\n", *req.RuleName, err)
			return err
		}
		req.Tags = v.Tags
		if err = z.ruleImporter.CreateRule(req); err != nil {
			fmt.Printf("rule name: %v import failed, err: %v\n", *req.RuleName, err)
			return err
		}
		fmt.Printf("rule name: %v import success!\n", *req.RuleName)
	}

	return nil
}

func (z *ZoneCopyManager) convertRules(old []*teo.Rule) ([]*teo.Rule, error) {
	if len(old) != 1 {
		return nil, fmt.Errorf("rule abnormal format")
	}
	// 第一层的if中condition需要转换域名相关条件，action暂未发现需要转换的
	// 第二层的if中condition需要转换域名相关条件，action需转换源站GroupId

	// 第一层if的替换
	z.convertConditions(old[0].Conditions)
	// 第二层if的判断
	subRule := old[0].SubRules
	if len(subRule) == 0 {
		return old, nil
	}
	for _, v1 := range subRule {
		for _, v2 := range v1.Rules {
			z.convertConditions(v2.Conditions)
			if err := z.convertActions(v2.Actions); err != nil {
				return nil, err
			}
		}
	}

	return old, nil
}

func (z *ZoneCopyManager) convertConditions(conds []*teo.RuleAndConditions) error {
	for _, v1 := range conds {
		for _, v2 := range v1.Conditions {
			// 第一层if中condition的域名进行替换
			if *v2.Target == "host" {
				for k, _ := range v2.Values {
					newDomainName := z.getNewDomainName(*v2.Values[k])
					v2.Values[k] = common.StringPtr(newDomainName)
				}
			}
		}
	}
	return nil
}

func (z *ZoneCopyManager) convertActions(actions []*teo.Action) error {
	for i, _ := range actions {
		// 修改源站GroupId
		if actions[i].NormalAction != nil && *(actions[i].NormalAction.Action) == "Origin" {
			for _, v := range actions[i].NormalAction.Parameters {
				if *v.Name == "OriginGroupId" {
					oldGroupId := v.Values[0]
					nw, err := z.getNewGroupId(*oldGroupId)
					if err != nil {
						return err
					}
					v.Values[0] = common.StringPtr(nw)
				}
			}
		}
	}
	return nil
}

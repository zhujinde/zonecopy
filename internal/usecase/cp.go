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
	config              *entity.ZoneCopyConfig
	originImporter      *repository.OriginManager
	domainImporter      *repository.DomainManager
	ruleImporter        *repository.RuleEngineManager
	zoneSettingImporter *repository.ZoneSettingManager

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
		config:              c,
		originImporter:      repository.NewOriginManager(c.Account),
		domainImporter:      repository.NewDomainManager(c.Account),
		ruleImporter:        repository.NewRuleEngineManager(c.Account),
		zoneSettingImporter: repository.NewZoneSettingManager(c.Account),

		isOriginInit:   false,
		templateOrigin: make(map[string]string),
		targetOrigin:   make(map[string]string),
	}
}

// ImportOrigin 源站导入。
func (z *ZoneCopyManager) ImportOrigin() error {
	oldGroups, err := z.originImporter.DescribeOriginGroupList(z.config.TemplateZoneId)
	if err != nil {
		log.Printf("zone id: %v describe origin group failed, err: %v\n", z.config.TemplateZoneId, err)
		return err
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
			log.Printf("origin：%v import failed, err: %v\n", req.OriginGroupName, err)
			return err
		}
	}
	return nil
}

// ImportDomains 域名导入。
func (z *ZoneCopyManager) ImportDomains() error {
	oldDomains, err := z.domainImporter.DescribeDomainListDetail(z.config.TemplateZoneId)
	if err != nil {
		log.Printf("zone id: %v describe domain list failed, err: %v\n", z.config.TemplateZoneId, err)
		return err
	}
	// TODO：验证对象存储源站是否正常
	for _, v := range oldDomains {
		newDomainName := z.getNewName(*v.DomainName)
		req := teo.NewCreateAccelerationDomainRequest()
		req.ZoneId = common.StringPtr(z.config.TargetZoneId)
		req.DomainName = common.StringPtr(newDomainName)
		req.OriginInfo, err = z.converDomainOrigin(v.OriginDetail)
		if err != nil {
			log.Printf("domain：%v -> %v convert config failed， err: %v\n", *v.DomainName, *req.DomainName, err)
			return err
		}
		if err = z.domainImporter.CreateDomain(req); err != nil {
			log.Printf("domain：%v -> %v import failed， err: %v\n", *v.DomainName, *req.DomainName, err)
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
			log.Printf("zone id: %v describe origin group failed, err: %v\n", z.config.TemplateZoneId, err)
			return "", err
		}
		for _, v := range oldGroups {
			z.templateOrigin[*v.OriginGroupId] = *v.OriginGroupName
		}
		newGroups, err := z.originImporter.DescribeOriginGroupList(z.config.TargetZoneId)
		if err != nil {
			log.Printf("zone id: %v describe origin group failed, err: %v\n", z.config.TargetZoneId, err)
			return "", err
		}
		for _, v := range newGroups {
			z.targetOrigin[*v.OriginGroupName] = *v.OriginGroupId
		}
		z.isOriginInit = true
	}
	name, ok := z.templateOrigin[old]
	if !ok {
		log.Printf("zoneId: %v not find old origin name: %v", z.config.TargetZoneId, old)
		return "", fmt.Errorf("not find old origin name")
	}
	id, ok := z.targetOrigin[name]
	if !ok {
		log.Printf("zoneId: %v not find new origin id: %v", z.config.TargetZoneId, old)
		return "", fmt.Errorf("not find new origin id")
	}
	return id, nil
}

// getNewDomainName 旧站点域名转换为新站点域名。
func (z *ZoneCopyManager) getNewName(old string) string {
	// TODO: 域名转换方式会有小概率bug
	return strings.Replace(old, z.config.TemplateZone, z.config.TargetZone, 1)
}

// ImportRuleEngineRules 规则引擎中规则的导入。
func (z *ZoneCopyManager) ImportRuleEngineRules() error {
	oldRules, err := z.ruleImporter.DescribeRuleList(z.config.TemplateZoneId)
	if err != nil {
		log.Printf("zone id: %v describe rule list failed, err: %v\n", z.config.TemplateZoneId, err)
		return err
	}
	// 逆序导入，最终展示保持和原站点一致
	l := len(oldRules)
	for i := l - 1; i >= 0; i-- {
		v := oldRules[i]
		// 规则名称如包含域名也进行一次替换
		newRuleName := z.getNewName(*v.RuleName)
		req := teo.NewCreateRuleRequest()
		req.ZoneId = common.StringPtr(z.config.TargetZoneId)
		req.RuleName = common.StringPtr(newRuleName)
		req.Status = common.StringPtr("enable")
		req.Rules, err = z.convertRules(v.Rules)
		if err != nil {
			log.Printf("rule name: %v convert config failed, err: %v\n", *req.RuleName, err)
			return err
		}
		req.Tags = v.Tags
		if err = z.ruleImporter.CreateRule(req); err != nil {
			log.Printf("rule name: %v import failed, err: %v\n", *req.RuleName, err)
			return err
		}
		log.Printf("rule name: %v import success!\n", *req.RuleName)
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
	if err := z.convertConditions(old[0].Conditions); err != nil {
		return nil, err
	}
	// 第二层if的判断
	subRule := old[0].SubRules
	if len(subRule) == 0 {
		return old, nil
	}
	for _, v1 := range subRule {
		for _, v2 := range v1.Rules {
			if err := z.convertConditions(v2.Conditions); err != nil {
				return nil, err
			}
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
					newDomainName := z.getNewName(*v2.Values[k])
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

// ImportZoneSetting 导入全局站点配置。
func (z *ZoneCopyManager) ImportZoneSetting() error {
	sets, err := z.zoneSettingImporter.DescribeZoneSetting(z.config.TemplateZoneId)
	if err != nil {
		log.Printf("zone id: %v describe zone setting failed, err: %v\n", z.config.TemplateZoneId, err)
		return err
	}
	req := teo.NewModifyZoneSettingRequest()
	req.ZoneId = common.StringPtr(z.config.TargetZoneId)
	req.CacheConfig = sets.CacheConfig
	req.CacheKey = sets.CacheKey
	req.MaxAge = sets.MaxAge
	req.OfflineCache = sets.OfflineCache
	req.Quic = sets.Quic
	req.PostMaxSize = sets.PostMaxSize
	req.Compression = sets.Compression
	req.UpstreamHttp2 = sets.UpstreamHttp2
	req.ForceRedirect = sets.ForceRedirect
	req.Https = sets.Https
	req.Origin = sets.Origin
	req.SmartRouting = sets.SmartRouting
	req.WebSocket = sets.WebSocket
	req.ClientIpHeader = sets.ClientIpHeader
	req.CachePrefresh = sets.CachePrefresh
	req.Ipv6 = sets.Ipv6
	req.ClientIpCountry = sets.ClientIpCountry
	req.Grpc = sets.Grpc
	// TODO: 媒体处理的配置当前版本接口不支持，无法拷贝
	// req.ImageOptimize = sets.ImageOptimize
	if err = z.zoneSettingImporter.ModifyZoneSetting(req); err != nil {
		log.Printf("zone id: %v import zone setting failed, err: %v\n", *req.ZoneId, err)
		return err
	}
	return nil
}

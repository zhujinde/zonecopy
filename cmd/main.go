package main

import (
	"flag"
	"fmt"

	"zonecopy/internal/domain/entity"
	"zonecopy/internal/usecase"
)

func main() {
	defer func() {
		if r := recover(); r != any(nil) {
			fmt.Println("[panic]", r)
		}
	}()

	var configPath, module string
	flag.StringVar(&module, "module", "", "导入指定模块配置 \norigin: 源站组 \ndomain: 域名管理 \nzonesetting: 站点加速配置 \nrule: 规则引擎 \nall: 全部模块")
	flag.StringVar(&configPath, "config", "./config/cp.yaml", "配置文件路径")
	//解析参数
	flag.Parse()
	var modules []FuncModule
	switch module {
	case "origin":
		modules = append(modules, moduleOrigin)
	case "domain":
		modules = append(modules, moduleDomain)
	case "zonesetting":
		modules = append(modules, moduleZoneSetting)
	case "rule":
		modules = append(modules, moduleRule)
	case "all":
		modules = []FuncModule{moduleOrigin, moduleDomain, moduleZoneSetting, moduleRule}
	default:
		panic(any("unsupported module!"))
	}

	c := entity.InitZoneCopyConfig(configPath)
	z := usecase.NewZoneCopyManager(c)
	for i, _ := range modules {
		modules[i](z)
	}
}

type FuncModule func(z *usecase.ZoneCopyManager)

var (
	moduleOrigin FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportOrigin(); err != nil {
			fmt.Printf("[Error] origin group import failed，err: %v\n", err)
		} else {
			fmt.Println("====> origin group import success!")
		}
	}
	moduleDomain FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportDomains(); err != nil {
			fmt.Printf("[Error] domain import failed，err: %v\n", err)
		} else {
			fmt.Println("====> domain import success!")
		}
	}
	moduleZoneSetting FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportZoneSetting(); err != nil {
			fmt.Printf("[Error] zone setting import failed，err: %v\n", err)
		} else {
			fmt.Println("====> zone setting import success!")
		}
	}
	moduleRule FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportRuleEngineRules(); err != nil {
			fmt.Printf("[Error] rule engine import failed，err: %v\n", err)
		} else {
			fmt.Println("====> rule engine import success!")
		}
	}
)

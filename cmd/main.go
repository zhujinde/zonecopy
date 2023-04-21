package main

import (
	"flag"
	"fmt"
	"log"

	"zonecopy/internal/domain/entity"
	"zonecopy/internal/usecase"
)

func main() {
	var module string
	flag.StringVar(&module, "module", "", "导入指定模块配置 \norigin: 源站组 \ndomain: 域名管理 \nrule: 规则引擎 \nall: 全部模块")
	//解析参数
	flag.Parse()
	var modules []FuncModule
	switch module {
	case "origin":
		modules = append(modules, moduleOrigin)
	case "domain":
		modules = append(modules, moduleDomain)
	case "rule":
		modules = append(modules, moduleRule)
	case "all":
		modules = []FuncModule{moduleOrigin, moduleDomain, moduleRule}
	default:
		panic(any("unsupported module!"))
	}

	c := entity.InitZoneCopyConfig(defaultConfigPath)
	z := usecase.NewZoneCopyManager(c)
	for i, _ := range modules {
		modules[i](z)
	}
}

type FuncModule func(z *usecase.ZoneCopyManager)

var (
	moduleOrigin FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportOrigin(); err != nil {
			fmt.Printf("====> origin group import failed，err: %v\n", err)
		} else {
			fmt.Println("### origin group import success!")
		}
	}
	moduleDomain FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportDomains(); err != nil {
			fmt.Printf("====> domain import failed，err: %v\n", err)
		} else {
			fmt.Println("### domain import success!")
		}
	}
	moduleRule FuncModule = func(z *usecase.ZoneCopyManager) {
		if err := z.ImportRuleEngineRules(); err != nil {
			log.Printf("====> rule engine import failed，err: %v\n", err)
		} else {
			log.Printf("### rule engine import success!")
		}
	}
	defaultConfigPath = "./config/cp.yaml"
)

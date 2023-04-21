package entity

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"zonecopy/pkg/utils"
)

// AccountBaseInfo 账号信息。
type AccountBaseInfo struct {
	SecretId  string `yaml:"secret_id" validate:"required"`
	SecretKey string `yaml:"secret_key" validate:"required"`
	EndPoint  string `yaml:"end_point" validate:"required"`
	Region    string `yaml:"region" validate:"required"`
}

// ZoneCopyConfig 初始化配置
type ZoneCopyConfig struct {
	LogPath        string           `yaml:"log_path" validate:"required"`
	Account        *AccountBaseInfo `yaml:"account" validate:"required"`
	TemplateZone   string           `yaml:"template_zone" validate:"required"`
	TemplateZoneId string           `yaml:"template_zone_id" validate:"required"`
	TargetZone     string           `yaml:"target_zone" validate:"required"`
	TargetZoneId   string           `yaml:"target_zone_id" validate:"required"`
	//Origin         map[string]string `yaml:"origin"`      // 源站组Name -> GroupId
	//RuleEngine     map[string]string `yaml:"rule_engine"` // 规则引擎 host -> RuleId
}

func InitZoneCopyConfig(configPath string) *ZoneCopyConfig {
	c := &ZoneCopyConfig{}
	if err := utils.PraseConfig(configPath, c); err != nil {
		panic(any(err))
	}
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		panic(any(err))
	}
	//if c.Origin == nil {
	//	c.Origin = make(map[string]string)
	//}
	//if c.RuleEngine == nil {
	//	c.RuleEngine = make(map[string]string)
	//}
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	logFile, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	boby, _ := json.Marshal(c)
	log.Printf("###### Init: %#v\n", string(boby))
	if err != nil {
		panic(any(err))
	}
	log.SetOutput(logFile)

	return c
}

// ZoneCopyConfigExport 初始化配置导出。
func ZoneCopyConfigExport(path string, z *ZoneCopyConfig) {
	err := utils.GenerateConfig(path, z)
	if err != nil {
		fmt.Errorf("utils.GenerateConfig failed, err: %v", err)
	}
}

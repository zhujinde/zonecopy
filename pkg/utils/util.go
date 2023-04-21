package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zclconf/go-cty/cty"
	yaml "gopkg.in/yaml.v2"
)

var (
	// OperatorMap akamai操作符映射Eo操作符。
	OperatorMap = map[string]string{
		"IS":                    "equal",
		"IS_NOT":                "notequal",
		"MATCHES_ONE_OF":        "exist",
		"DOES_NOT_MATCH_ONE_OF": "notexist",
		"IS_ONE_OF":             "exist",
		"IS_NOT_ONE_OF":         "notexist",
		"EXISTS":                "exist",
	}

	TimeMap = map[byte]int64{
		's': 1,
		'm': 60,
		'h': 3600,
		'd': 86400,
	}

	DefaultConfigFileType = ".json"
	configFileType        string
	configFiles           []string
)

// CtyStrList 生成cty字符串列表。
func CtyStrList(list []string) []cty.Value {
	var rsp []cty.Value
	for _, v := range list {
		rsp = append(rsp, cty.StringVal(v))
	}

	return rsp
}

// IsParaVariable 判断参数是否为变量。
func IsParaVariable(s string) bool {
	b := []byte(s)
	len := len(b)

	return len > 0 && b[0] == '{' && b[len-1] == '}'
}

// TranTimeToSec 时间字符串转换成秒数字符串。
func TranTimeToSec(ttl string) (string, error) {
	if len(ttl) < 2 {
		return "", fmt.Errorf("error ttl")
	}
	t := []byte(ttl)
	tLen := len(t)
	tNum := t[:tLen-1]
	tUnit := t[tLen-1]

	num, err := strconv.ParseInt(string(tNum), 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse time num failed")
	}
	unit, ok := TimeMap[tUnit]
	if !ok {
		return "", fmt.Errorf("parse time unit failed")
	}
	return strconv.FormatInt(num*unit, 10), nil
}

// Listfunc ...
func Listfunc(path string, f os.FileInfo, err error) error {
	ok := strings.HasSuffix(path, configFileType)
	if ok {
		configFiles = append(configFiles, path)
	}
	return nil
}

// GetFileList 扫描指定路径下文件列表。
func GetFileList(path string, fileType string) ([]string, error) {
	configFileType = DefaultConfigFileType
	if fileType != "" {
		configFileType = fileType
	}
	configFiles = configFiles[0:0]
	if err := filepath.Walk(path, Listfunc); err != nil {
		return nil, err
	}
	return configFiles, nil
}

// PraseConfig 从本地解析yaml配置。
func PraseConfig(path string, dest interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(f)
	return decoder.Decode(dest)
}

// GenerateConfig 导出yaml配置。
func GenerateConfig(path string, dest interface{}) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(f)
	return encoder.Encode(dest)
}

// GenTencentOutFileName 生成输出文件路径。
func GenTencentOutFileName(filename string, outputDir string) (string, string, error) {
	strs := strings.Split(filename, "/")
	l := len(strs)
	if l < 1 {
		return "", "", fmt.Errorf("invalid filename")
	}
	abs := strs[l-1]
	tfName := outputDir + "/" + abs + ".tf"
	failName := outputDir + "/" + abs + ".fail.json"
	return tfName, failName, nil
}

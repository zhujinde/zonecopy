package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	DefaultConfigFileType = ".json"
	configFileType        string
	configFiles           []string
)

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

// GenTencentOutFileName 生成配置输出文件。
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

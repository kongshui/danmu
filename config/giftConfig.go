package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 配置文件结构体
type CfgConfig struct {
	Type   map[string]string `json:"type"`
	Fields []map[string]any  `json:"fields"`
}

// 读取配置文件
func ReadCfgConfig(filePath string) (*CfgConfig, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ReadCfgConfig 打开文件失败：%v", err)
	}
	defer file.Close()
	var (
		typeList []string // 配置类型列表
		fileList []string // 配置内容列表
	)
	// 礼物配置实例
	cfgConfig := CfgConfig{
		Type: make(map[string]string),
	}
	lineNum := 0
	// 解析文件
	scanner := bufio.NewScanner(file)
	// 遍历文件内容
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 {
			continue
		}
		line := scanner.Text()
		re := regexp.MustCompile(`[\t\s]+`)
		fields := re.Split(line, -1)
		if lineNum == 2 {
			typeList = fields
		}
		if lineNum == 3 {
			fileList = fields
			if len(typeList) != len(fileList) {
				return nil, fmt.Errorf("ReadCfgConfig 配置文件格式错误：第%d行，字段数与配置类型数不一致", lineNum)
			}
		}
		if len(fields) != len(typeList) {
			return nil, fmt.Errorf("ReadCfgConfig 配置文件格式错误：第%d行，字段数与配置类型数不一致", lineNum)
		}
		if lineNum > 3 {
			data := make(map[string]any)
			for i, field := range fields {
				data[fileList[i]] = field
				switch strings.ToLower(typeList[i]) {
				case "int":
					intVal, err := strconv.Atoi(field)
					if err != nil {
						return nil, fmt.Errorf("ReadCfgConfig 转换int失败：%v", err)
					}
					data[fileList[i]] = intVal
				case "float":
					floatVal, err := strconv.ParseFloat(field, 64)
					if err != nil {
						return nil, fmt.Errorf("ReadCfgConfig 转换float失败：%v", err)
					}
					data[fileList[i]] = floatVal
				case "bool":
					boolVal, err := strconv.ParseBool(field)
					if err != nil {
						return nil, fmt.Errorf("ReadCfgConfig 转换bool失败：%v", err)
					}
					data[fileList[i]] = boolVal
				default:
					data[fileList[i]] = field
				}
			}
			cfgConfig.Fields = append(cfgConfig.Fields, data)
		}
		for i, field := range fileList {
			cfgConfig.Type[field] = typeList[i]
		}
		// 检查扫描过程中是否出错（如IO错误）
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("ReadCfgConfig 逐行读取失败：%v", err)
		}
	}
	return &cfgConfig, nil
}

package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveToJSON 序列化任意结构体到 JSON 文件
func SaveToJSON(data interface{}, filePath string) error {
	// 将数据序列化为 JSON 字节
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// 创建或打开文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 将 JSON 数据写入文件
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("json文件保存成功，文件路径为：", filePath)

	return nil
}

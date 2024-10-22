package conf

import (
	"encoding/json"
	"os"
)

// Config 定义公共配置结构体
type Config struct {
	Token    string `json:"token"`
	MaxRetry int32  `json:"maxRetry"`
	Cookies  string `json:"cookies"`
}

// AppConfig 全局变量，用于存储加载的配置
var AppConfig Config

// LoadConfig 从JSON文件加载配置
func LoadConfig(filePath string) error {
	configFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}

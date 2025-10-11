package conf

import (
	"encoding/json"
	"os"
	"xizexcample/internal/pkg/logger"
)

// Config 保存应用程序的所有配置
type Config struct {
	ServerHost string `json:"server_host"`
	ServerPort int    `json:"server_port"`
	Timeout    int    `json:"timeout"`
}

// AppConfig 是全局应用程序配置
var AppConfig *Config

func init() {
	// 为 AppConfig 提供默认值
	AppConfig = &Config{
		ServerHost: "0.0.0.0",
		ServerPort: 8999,
		Timeout:    5,
	}
	LoadConfig("conf/zinx.json")
}

// LoadConfig 从文件中加载配置
func LoadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to open config file: %v", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(AppConfig)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to decode config file: %v", err)
		return
	}
	logger.InfoLogger.Printf("Configuration loaded from %s", path)
}

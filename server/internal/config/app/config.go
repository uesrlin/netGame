package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	FileName string `yaml:"fileName"`
	FilePath string `yaml:"filePath"`
}

type Config struct {
	LogConfig   LogConfig `yaml:"log"`
	Debug       bool      `yaml:"debug"`
	Host        string    `yaml:"host"`
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	ListenPoint string    `yaml:"listenPoints"`
}

// ConfigMap /**
type ConfigMap struct {
	App Config `yaml:"app"`
}

func InitAppConfigMap(configData []byte) *ConfigMap {
	value := ConfigMap{
		App: Config{},
	}
	err := yaml.Unmarshal(configData, &value)
	if err != nil {
		panic(fmt.Sprintf("解析App配置文件失败%s", err))
	}
	return &value
}

func (c *ConfigMap) GetLogFilePath() string {
	return c.App.LogConfig.FilePath
}

func (c *ConfigMap) GetLogFileName() string {
	return c.App.LogConfig.FileName
}

func (c *ConfigMap) GetHost() string {
	return c.App.Host
}

func (c *ConfigMap) GetDebug() bool {
	return c.App.Debug
}

func (c *ConfigMap) GetListenPoint() string {
	return c.App.ListenPoint
}

func (c *ConfigMap) Version() string {
	return c.App.Version
}

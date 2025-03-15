package app

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	FileName string `yaml:"fileName"`
	FilePath string `yaml:"filePath"`
}

type Config struct {
	LogConfig       LogConfig `yaml:"log"`
	Salt            string    `yaml:"salt"`
	Debug           bool      `yaml:"debug"`
	ListenPoint     string    `yaml:"listenPoint"`
	InternalService bool      `yaml:"internalService"`
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

func (c *ConfigMap) GetSalt() string {
	return c.App.Salt
}

func (c *ConfigMap) GetDebug() bool {
	return c.App.Debug
}

func (c *ConfigMap) GetListenPoint() string {
	return c.App.ListenPoint
}

func (c *ConfigMap) IsInternalService() bool {
	return c.App.InternalService
}

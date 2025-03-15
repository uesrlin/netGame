package redis

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

/**
 * @Description
 * @Date 2025/3/15 19:25
 **/

type Config struct {
	Host        string `yaml:"host"`
	MaxPoolSize uint64 `yaml:"maxPoolSize"`
	MinPoolSize uint64 `yaml:"minPoolSize"`
	Passwd      string `yaml:"passwd"`
	Timeout     int    `yaml:"timeout"`
}

func (m *Config) ApplyURL() string {
	return m.Host
}

// ConfigMap /**
type ConfigMap struct {
	Dbs map[string]Config `yaml:"dbs"`
}

func (m *ConfigMap) GetDBConfigWithName(name string) *Config {
	value, ok := m.Dbs[name]
	if !ok {
		return nil
	}
	return &value
}

func InitDBConfigMap(configData []byte) *ConfigMap {
	value := ConfigMap{
		Dbs: make(map[string]Config, 10),
	}
	err := yaml.Unmarshal(configData, &value)
	if err != nil {
		panic(fmt.Sprintf("解析Redis配置文件失败%s", err))
	}
	return &value
}

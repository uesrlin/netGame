package mysql

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

/**
 * @Description
 * @Date 2025/3/15 19:25
 **/

type Config struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	DBName      string `yaml:"dbName"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
	MinPoolSize int    `yaml:"minPoolSize"`
	User        string `yaml:"user"`
	Passwd      string `yaml:"passwd"`
	Timeout     int    `yaml:"timeout"`
}

func (m *Config) ApplyURL() string {
	args := strings.Builder{}
	args.WriteString(m.User)
	args.WriteString(":")
	args.WriteString(m.Passwd)
	args.WriteString("@tcp(")
	args.WriteString(m.Host)
	args.WriteString(":")
	args.WriteString(m.Port)
	args.WriteString(")/")
	args.WriteString(m.DBName)
	args.WriteString("?")
	args.WriteString("parseTime=true&loc=Local")
	return args.String()
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
		panic(fmt.Sprintf("解析Mysql配置文件失败%s", err))
	}
	return &value
}

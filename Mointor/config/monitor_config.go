package configs

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// 全局配置
var Config *ServerConfig

type ServerConfig struct {
	Time int64 `yaml:"time:"`
}

func LoadConfig(filePath string) (config *ServerConfig, err error) {
	config = &ServerConfig{
		Time: 60,
	}
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("LoadConfig(%s) failed! err: %v", filePath, err)
		return nil, err
	}
	err = yaml.Unmarshal(buf, Config)
	if err != nil {
		log.Fatalf("Unmarshal yaml config failed! err: %v", filePath, err)
		return nil, err
	}
	return Config, nil

}

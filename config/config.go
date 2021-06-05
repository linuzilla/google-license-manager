package config

import (
	"github.com/linuzilla/go-logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Name         string `yaml:"name"`
	LogLevel     string `yaml:"log-level"`
	DatabaseFile string `yaml:"database-file"`
	GoogleCfg    `yaml:"google"`
}

func FromFile(fileName string) (*Config, error) {
	yamlContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return FromData(yamlContent)
}

func FromData(yamlContent []byte) (*Config, error) {
	conf := new(Config)
	err := yaml.Unmarshal(yamlContent, conf)
	if err != nil {
		logger.Error("Unmarshal: %v", err)
		return nil, err
	}
	return conf, nil
}

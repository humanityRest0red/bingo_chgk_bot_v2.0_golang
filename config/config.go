package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		Path string `yaml:"path"`
	} `yaml:"log"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

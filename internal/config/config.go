package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AppEnv   AppEnv   `yaml:"app-env"`
	DBConfig DBConfig `yaml:"db-config"`
}

type AppEnv struct {
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
}

type DBConfig struct {
	DatabaseUrl string `yaml:"database-url"`
}

func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open config: %w", err))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(fmt.Errorf("failed to decode config.yaml: %w", err))
	}

	return &cfg
}

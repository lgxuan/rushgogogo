package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		ThreadCount int `yaml:"threadCount"`
	} `yaml:"server"`
	Filters []Filter `yaml:"filters"`
}

type Filter struct {
	FilterName     string `yaml:"filterName"`
	FilterType     string `yaml:"filterType"`
	FilterRegex    string `yaml:"filterRegex"`
	FilterResource string `yaml:"filterResource"`
	FilterLevel    string `yaml:"filterLevel"`
	FilterEnabled  bool   `yaml:"filterEnabled"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func UpdateThreadCount(configPath string, threadCount int) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	config.Server.ThreadCount = threadCount

	return SaveConfig(config, configPath)
}

func GetThreadCount(configPath string) (int, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return 0, err
	}

	return config.Server.ThreadCount, nil
}

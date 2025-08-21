package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type filter struct {
	FilterName     string `yaml:"filterName,omitempty"`
	FilterType     string `yaml:"filterType,omitempty"`
	FilterRegex    string `yaml:"filterRegex,omitempty"`
	FilterResource string `yaml:"filterResource,omitempty"`
	FilterLevel    string `yaml:"filterLevel,omitempty"`
	FilterEnabled  bool   `yaml:"filterEnabled,omitempty"`
}
type Configuration struct {
	Filters []filter `yaml:"filters,omitempty"`
	Cert    cert     `yaml:"customCert,omitempty"`
}
type cert struct {
	Cert string `yaml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

func GetConfigurationFromYaml(fileName string) (Configuration, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Configuration{}, fmt.Errorf("config file read error: %w", err)
	}

	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, fmt.Errorf("config file parse error: %w", err)
	}
	return config, nil
}

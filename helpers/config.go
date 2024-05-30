package helpers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service struct {
		Port string `yaml:"port"`
	} `yaml:"service"`
	TargetLabels []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"targetLabels"`
	Tolerations []struct {
		Key      string `yaml:"key"`
		Operator string `yaml:"operator"`
		Value    string `yaml:"value"`
		Effect   string `yaml:"effect"`
	} `yaml:"tolerations"`
	NodeSelector struct {
		AgentPool string `yaml:"agentpool"`
	} `yaml:"nodeSelector"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}
	return &config, nil
}
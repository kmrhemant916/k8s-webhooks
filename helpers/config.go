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
	Patch struct {
		Tolerations struct {
			Enable bool `yaml:"enable"`
			Value  []struct {
				Key      string `yaml:"key"`
				Operator string `yaml:"operator"`
				Value    string `yaml:"value"`
				Effect   string `yaml:"effect"`
			} `yaml:"value"`
		} `yaml:"tolerations"`
		NodeSelector struct {
			Enable bool `yaml:"enable"`
			Value  struct {
				AgentPool string `yaml:"agentpool"`
			} `yaml:"value"`
		} `yaml:"nodeSelector"`
		ImagePullSecrets struct {
			Enable bool `yaml:"enable"`
			Value  []struct {
				Name string `yaml:"name"`
			} `yaml:"value"`
		} `yaml:"imagePullSecrets"`
	} `yaml:"patch"`
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
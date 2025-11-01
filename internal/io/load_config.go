package io

import (
	"mockium/pkg/ports"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

var LoadConfig = func(path string) (*ports.Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ports.Config{}
	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		return nil, err
	}

	return config, nil
}

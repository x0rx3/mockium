package io

import (
	"mockium/pkg/core"
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

var LoadConfig = func(path string) *core.Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	config := &core.Config{}
	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		panic(err)
	}

	return config
}

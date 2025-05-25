package model

import "os"

type SetResponse struct {
	SetStatus  int
	SetHeaders map[string]string
	SetBody    map[string]any
	SetFile    *os.File
}

type SetResponseTemplate struct {
	SetStatus  int               `yaml:"SetStatus" json:"SeStatus"`
	SetHeaders map[string]string `yaml:"SetHeaders" json:"SetHeaders"`
	SetBody    map[string]any    `yaml:"SetBody" json:"SetBody"`
	SetFile    string            `yaml:"SetFile" json:"SetFile"`
}

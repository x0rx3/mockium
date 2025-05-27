package model

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func (inst *SetResponseTemplate) UnmarshalJSON(data []byte) error {
	type Alias SetResponseTemplate
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(inst),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if inst.SetFile != "" && inst.SetBody != nil {
		return fmt.Errorf("cannot use parameter 'SetBody' with 'SetFile'")
	}

	return nil
}

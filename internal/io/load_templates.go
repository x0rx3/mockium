package io

import (
	"encoding/json"
	"fmt"
	"mockium/pkg/model"
	"os"
	"regexp"
)

var LoadTemplates = func(path string) ([]model.Template[[]model.HandleTemplate], error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(`.json`)
	if err != nil {
		return nil, err
	}

	templates := make([]model.Template[[]model.HandleTemplate], 0)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		template := model.Template[[]model.HandleTemplate]{}
		if re.MatchString(file.Name()) {
			f, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("%s, file: %s", err.Error(), file.Name())
			}

			if err := json.Unmarshal(f, &template); err != nil {
				return nil, fmt.Errorf("%s, file: %s", err.Error(), file.Name())
			}

			templates = append(templates, template)
		}
	}

	return templates, nil
}

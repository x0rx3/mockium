package builder

import (
	"encoding/json"
	"fmt"
	"gomock/internal/model"
	"os"
	"regexp"

	"go.uber.org/zap"
)

// NewTemplateBuilder creates a new TemplateBuilder instance.
// It initializes the builder with the provided logger and sets up the template building logic.
func NewTemplateBuilder(log *zap.Logger) *TemplateBuilder {
	return &TemplateBuilder{
		log: log,
	}
}

// TemplateBuilder is a struct that implements the TemplateBuilder interface.
// It is responsible for building templates based on the provided path.
// The struct contains the logger and provides methods to build templates.
// It reads JSON files from the specified path and unmarshals them into HandlerTamplate structs.
// The struct also validates the templates to ensure that they meet the required criteria.
type TemplateBuilder struct {
	log *zap.Logger
}

func (inst *TemplateBuilder) Build(path string) ([]model.Template, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(`.json`)
	if err != nil {
		return nil, err
	}

	templates := make([]model.Template, 0)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		template := model.Template{}
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

	if err := inst.validate(templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// validate checks if the provided templates meet the required criteria.
// It ensures that the templates do not contain conflicting parameters.
// For example, it checks if both SetBody and SetFile are used in the same template.
// If any conflicts are found, it returns an error.
// If the templates are valid, it returns nil.
func (inst *TemplateBuilder) validate(templates []model.Template) error {
	for _, template := range templates {
		for _, handle := range template.Handle {
			if handle.SetResponseTemplate.SetBody != nil && handle.SetResponseTemplate.SetFile != "" {
				return fmt.Errorf("cannot use parameter 'SetBody' with 'SetFile'")
			}
		}
	}
	return nil
}

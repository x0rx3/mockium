package builder

import (
	"encoding/json"
	"fmt"
	"mockium/internal/model"
	"os"
	"regexp"

	"go.uber.org/zap"
)

// TemplateBuilder is responsible for loading and validating template definitions
// from JSON files in a specified directory.
type TemplateBuilder struct {
	log *zap.Logger // Logger for validation and loading diagnostics (currently unused).
}

// NewTemplateBuilder creates a new instance of TemplateBuilder.
//
// Parameters:
//   - log: zap logger used for debug or error logging.
//
// Returns a pointer to a TemplateBuilder.
func NewTemplateBuilder(log *zap.Logger) *TemplateBuilder {
	return &TemplateBuilder{
		log: log,
	}
}

// Build reads all JSON template files from the given directory path, unmarshals them,
// and validates the resulting templates.
//
// Parameters:
//   - path: directory path where template JSON files are located.
//
// Returns a slice of model.Template and an error if reading or validation fails.
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

// validate performs structural validation of templates including:
//   - setting default HTTP method if not specified
//   - ensuring only one of SetBody or SetFile is used in a response
//   - checking for valid HTTP methods
//
// Parameters:
//   - templates: the slice of templates to validate.
//
// Returns an error if validation fails.
func (inst *TemplateBuilder) validate(templates []model.Template) error {
	for _, template := range templates {
		for _, handle := range template.Handle {

			if handle.MatchRequestTemplate.MustMethod == "" {
				handle.MatchRequestTemplate.MustMethod = model.DEFAULTMETHOD
			} else {
				if err := inst.checkMethod(handle.MatchRequestTemplate.MustMethod); err != nil {
					return err
				}
			}

			if handle.SetResponseTemplate.SetBody != nil && handle.SetResponseTemplate.SetFile != "" {
				return fmt.Errorf("cannot use parameter 'SetBody' with 'SetFile'")
			}
		}
	}
	return nil
}

// checkMethod verifies that the provided HTTP method is supported.
//
// Parameters:
//   - metod: the HTTP method to validate.
//
// Returns an error if the method is not recognized.
func (inst *TemplateBuilder) checkMethod(metod model.Method) error {
	switch metod {
	case model.GET, model.POST, model.DELETE, model.PATCH, model.PUT:
		return nil
	default:
		return fmt.Errorf("unexpected method")
	}
}

package httpval

import (
	"mockium/internal/defs"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func New() core.Validator {
	return &httpValidator{}
}

type httpValidator struct{}

func (inst *httpValidator) Validate(template model.Template[[]model.HandleTemplate]) error {
	for _, handle := range template.Handle {
		switch template.Type {
		case string(defs.HTTP):
			if err := inst.validateHandleTemplate(handle); err != nil {
				return err
			}
		default:
			return defs.ErrorUnxpectedHandleType
		}
		inst.setDefaultValues(&handle)
	}
	return nil
}

func (inst *httpValidator) validateHandleTemplate(handle model.HandleTemplate) error {
	if err := inst.validateResponse(handle.SetResponseTemplate); err != nil {
		return err
	}

	if err := inst.validateRequest(handle.MatchRequestTemplate); err != nil {
		return err
	}
	return nil
}

func (inst *httpValidator) validateResponse(response model.SetResponseTemplate) error {
	if response.SetBody != nil && response.SetFile != "" {
		return defs.ErrorSetBodyWithSetFile
	}

	return nil
}

func (inst *httpValidator) validateRequest(request model.MatchRequestTemplate) error {
	if err := inst.validateMethod(request.MustMethod); err != nil {
		return err
	}

	return nil
}

func (inst *httpValidator) validateMethod(metod string) error {
	switch metod {
	case string(defs.GET),
		string(defs.POST),
		string(defs.PUT),
		string(defs.DELETE),
		string(defs.PATCH),
		string(defs.EMPTYMETHOD): // default method
		return nil
	default:
		return defs.ErrorUnxpectedHTTPMethod
	}
}

func (inst *httpValidator) setDefaultValues(handle *model.HandleTemplate) {
	if handle.MatchRequestTemplate.MustMethod == "" {
		handle.MatchRequestTemplate.MustMethod = string(defs.DEFAULTMETHOD)
	}

	if handle.SetResponseTemplate.SetStatus == 0 {
		handle.SetResponseTemplate.SetStatus = 200
	}
}

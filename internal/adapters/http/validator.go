package http

import (
	"mockium/internal/adapters/common"
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

func NewValidator() ports.Validator {
	return &httpValidator{}
}

type httpValidator struct{}

func (inst *httpValidator) Validate(template model.Template[[]model.HandleTemplate]) error {
	for _, handle := range template.Handle {
		switch template.Type {
		case string(common.HTTP):
			if err := inst.validateHandleTemplate(handle); err != nil {
				return err
			}
		default:
			return common.ErrorUnxpectedHandleType
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
		return common.ErrorSetBodyWithSetFile
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
	case string(common.GET),
		string(common.POST),
		string(common.PUT),
		string(common.DELETE),
		string(common.PATCH),
		string(common.EMPTYMETHOD): // default method
		return nil
	default:
		return common.ErrorUnxpectedHTTPMethod
	}
}

func (inst *httpValidator) setDefaultValues(handle *model.HandleTemplate) {
	if handle.MatchRequestTemplate.MustMethod == "" {
		handle.MatchRequestTemplate.MustMethod = string(common.DEFAULTMETHOD)
	}

	if handle.SetResponseTemplate.SetStatus == 0 {
		handle.SetResponseTemplate.SetStatus = 200
	}
}

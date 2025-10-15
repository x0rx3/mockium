package grpcval

import (
	"mockium/internal/defs"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func New() core.Validator {
	return &grpcValidator{}
}

type grpcValidator struct{}

func (inst *grpcValidator) Validate(template model.Template[[]model.HandleTemplate]) error {
	for _, handle := range template.Handle {
		if template.Type != string(defs.GRPC) {
			return defs.ErrorUnxpectedHandleType
		}

		if err := inst.validateHandleTemplate(handle); err != nil {
			return err
		}

		inst.setDefaultValues(&handle)
	}
	return nil
}

func (inst *grpcValidator) validateHandleTemplate(handle model.HandleTemplate) error {
	return nil
}

func (inst *grpcValidator) setDefaultValues(handle *model.HandleTemplate) {}

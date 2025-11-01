package grpc

import (
	"mockium/internal/adapters/common"
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

func New() ports.Validator {
	return &grpcValidator{}
}

type grpcValidator struct{}

func (inst *grpcValidator) Validate(template model.Template[[]model.HandleTemplate]) error {
	for _, handle := range template.Handle {
		if template.Type != string(common.GRPC) {
			return common.ErrorUnxpectedHandleType
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

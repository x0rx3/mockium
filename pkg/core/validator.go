package core

import "mockium/pkg/core/model"

type Validator interface {
	Validate(template model.Template[[]model.HandleTemplate]) error
}

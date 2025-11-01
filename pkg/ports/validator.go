package ports

import "mockium/pkg/model"

type Validator interface {
	Validate(template model.Template[[]model.HandleTemplate]) error
}

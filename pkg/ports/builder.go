package ports

import (
	"mockium/pkg/model"
)

type ResponseBuilder[T Request] interface {
	Build(T) (*model.SetResponse, error)
}

type HandlerBuilder[R Request, T Handler[R]] interface {
	Build(logger Logger, path string, handles []model.Handle) T
}

type ServerBuilder interface {
	Build(logger Logger, template []model.Template[[]model.Handle]) Server
}

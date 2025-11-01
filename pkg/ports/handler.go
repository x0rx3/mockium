package ports

import (
	"mockium/pkg/model"
)

type Handler[T Request] interface {
	Handle(T) (*model.SetResponse, error)
}

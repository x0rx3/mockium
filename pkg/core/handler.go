package core

import (
	"mockium/pkg/core/model"
)

type Handler[T Request] interface {
	Handle(T) (*model.SetResponse, error)
}

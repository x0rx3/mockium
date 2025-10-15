package core

import (
	"mockium/pkg/core/model"
)

type Service[T Request] interface {
	Handle(T) (*model.SetResponse, error)
}

type ResMtchMap[T Request] map[ResponseBuilder[T]][]Matcher[T]

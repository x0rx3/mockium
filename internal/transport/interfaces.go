package transport

import (
	"gomock/internal/model"
	"gomock/internal/transport/method"
	"net/http"
)

type ResponseBuilder interface {
	Build(req *http.Request) (*model.SetResponse, error)
}

type RequestMatcher interface {
	Match(req *http.Request) bool
}

type Router interface {
	Path() string
	Method() method.Method
	Handler() http.Handler
}

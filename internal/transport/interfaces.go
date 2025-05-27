package transport

import (
	"mockium/internal/model"
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
	Handlers() map[model.Method]http.Handler
	Handler(model.Method) http.Handler
}

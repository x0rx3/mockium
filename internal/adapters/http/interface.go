package http

import (
	"context"
	"mockium/pkg/model"
	"mockium/pkg/ports"
	"net/http"
)

type HTTPRequest interface {
	ports.Request
	Method() string
	Path() string
	WithContext(context.Context)
	Raw() *http.Request
}

type HTTPHandler[T HTTPRequest] interface {
	ports.Handler[T]
	Path() string
	SupportedMethod() []string
}

type Responder interface {
	Write(w http.ResponseWriter, r *http.Request, resp model.SetResponse)
	Error(w http.ResponseWriter, status int, msg string)
}

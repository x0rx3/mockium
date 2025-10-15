package httpgw

import (
	"context"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
	"net/http"
)

type HTTPRequest interface {
	core.Request
	Method() string
	Path() string
	WithContext(context.Context)
	Raw() *http.Request
}

type HTTPHandler[T HTTPRequest] interface {
	core.Handler[T]
	Path() string
	SupportedMethod() []string
}

type Responder interface {
	Write(w http.ResponseWriter, r *http.Request, resp model.SetResponse)
	Error(w http.ResponseWriter, status int, msg string)
}

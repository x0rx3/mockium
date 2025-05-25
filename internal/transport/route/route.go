package route

import (
	"gomock/internal/transport/method"
	"net/http"
)

// Route represents an HTTP route with a path, method, and handler.
// It implements the Router interface, which defines the methods for getting
// the path, method, and handler of the route.
// The Route struct is used to define a specific route in the HTTP server.
// It contains the path, method, and handler for the route.
// The path is the URL pattern that the route matches, the method is the HTTP method
// (e.g., GET, POST) that the route responds to, and the handler is the function
// that handles the request when the route is matched.
type Route struct {
	path    string
	method  method.Method
	handler http.Handler
}

// New creates a new Route instance with the specified path, method, and handler.
func New(path string, method method.Method, handler http.Handler) *Route {
	return &Route{
		path:    path,
		method:  method,
		handler: handler,
	}
}

// Path returns the path of the route.
func (inst *Route) Path() string { return inst.path }

// Method returns the HTTP method of the route.
func (inst *Route) Method() method.Method { return inst.method }

// Handler returns the HTTP handler for the route.
func (inst *Route) Handler() http.Handler { return inst.handler }

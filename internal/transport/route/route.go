package route

import (
	"mockium/internal/model"
	"net/http"
)

// New creates a new Route instance with the specified path and method handlers.
//
// Parameters:
//   - path: URL path pattern for the route (e.g., "/users/{id}")
//   - handlers: Map of HTTP methods to their corresponding handlers
//
// Returns:
//   - Pointer to a new Route instance
func New(path string, handlers map[model.Method]http.Handler) *Route {
	return &Route{
		path:     path,
		handlers: handlers,
	}
}

// Route represents an HTTP route configuration.
// It encapsulates:
// - The path pattern to match against incoming requests
// - A collection of handlers for different HTTP methods
//
// The struct implements the Router interface, providing access to:
// - The route path via Path()
// - All handlers via Handlers()
// - Specific handler by method via Handler()
type Route struct {
	path     string                        // URL path pattern
	handlers map[model.Method]http.Handler // Method-to-handler mappings
}

// Path returns the route's URL path pattern.
// This is used by the router to match incoming requests.
func (inst *Route) Path() string { return inst.path }

// Handlers returns all HTTP handlers configured for this route,
// indexed by HTTP method.
//
// Returns:
//   - Map of HTTP methods to their handlers
func (inst *Route) Handlers() map[model.Method]http.Handler { return inst.handlers }

// Handler returns the HTTP handler for a specific method.
//
// Parameters:
//   - method: HTTP method to get handler for
//
// Returns:
//   - Handler for the specified method, or nil if not found
func (inst *Route) Handler(method model.Method) http.Handler { return inst.handlers[method] }

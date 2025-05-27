package server

import (
	"mockium/internal/transport"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// New creates a new Server instance with the specified logger and routers.
//
// Parameters:
//   - log: Logger instance for logging server operations
//   - routes: Variadic list of routers to be registered with the server
//
// Returns:
//   - Pointer to a newly initialized Server instance
func New(log *zap.Logger, routes ...transport.Router) *Server {
	return &Server{
		log:    log,
		server: &http.Server{},
		routes: routes,
	}
}

// Server represents an HTTP server that manages multiple routers.
// It encapsulates:
// - A logger for recording server operations
// - The underlying http.Server instance
// - A collection of registered routers
type Server struct {
	log    *zap.Logger        // Logger for server operations
	server *http.Server       // Underlying HTTP server
	routes []transport.Router // Collection of registered routers
}

// Start initializes and runs the HTTP server on the specified address.
// It performs the following operations:
// 1. Configures the server address
// 2. Creates a new router using gorilla/mux
// 3. Registers all handlers from the configured routes
// 4. Starts listening for incoming requests
//
// Parameters:
//   - address: Network address to listen on (e.g., ":8080")
//
// Returns:
//   - error: Any error that occurs during server startup or operation
//
// Notes:
// - Defaults to GET method if no method is specified in the route
// - Logs each registered handler for debugging purposes
func (inst *Server) Start(address string) error {
	inst.server.Addr = address

	// Initialize the request router
	r := mux.NewRouter()

	// Register all routes and their handlers
	method := "GET"
	for _, route := range inst.routes {
		for m, hr := range route.Handlers() {
			// Use GET as default method if not specified
			if string(m) == "" {
				method = "GET"
			} else {
				method = string(m)
			}

			// Register the handler with the router
			r.HandleFunc(route.Path(), hr.ServeHTTP).Methods(method)

			// Log the registered handler
			inst.log.Info("added handler:",
				zap.String("path", route.Path()),
				zap.String("method", method))
		}
	}

	// Set the configured router as the server handler
	inst.server.Handler = r

	// Start the server
	inst.log.Info("start listen and serve",
		zap.String("address", address))

	return inst.server.ListenAndServe()
}

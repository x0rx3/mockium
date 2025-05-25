package server

import (
	"gomock/internal/transport"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// New creates a new Server instance with the given logger and routes.
// It initializes the server with the provided routes and sets up the logger.
func New(log *zap.Logger, routes ...transport.Router) *Server {
	return &Server{
		log:    log,
		server: &http.Server{},
		routes: routes,
	}
}

// Server is a simple HTTP server that uses the gorilla/mux router.
// It is designed to be used with the transport package to handle HTTP requests.
// It is not intended to be used as a standalone server.
type Server struct {
	log    *zap.Logger
	server *http.Server
	routes []transport.Router
}

// Start starts the HTTP server on the specified address.
// The server will listen on the specified address and handle incoming requests using the registered routes.
// It sets up the router and registers the routes with their corresponding handlers.
// It also logs the added handlers and starts listening for incoming requests.
// If an error occurs during startup, it returns the error.
func (inst *Server) Start(address string) error {
	inst.server.Addr = address

	r := mux.NewRouter()

	for _, hr := range inst.routes {
		r.HandleFunc(hr.Path(), hr.Handler().ServeHTTP).Methods(string(hr.Method()))
		inst.log.Info("added handler:", zap.String("path", hr.Path()), zap.String("method", string(hr.Method())))
	}
	inst.server.Handler = r

	inst.log.Info("start listen and serve", zap.String("address", address))

	return inst.server.ListenAndServe()
}

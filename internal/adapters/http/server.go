package http

import (
	"context"
	"mockium/internal/adapters/common"
	"mockium/pkg/ports"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewServer(log ports.Logger, address string, handlers ...HTTPHandler[HTTPRequest]) ports.Server {
	return &server{
		logger: log,
		server: &http.Server{
			Addr: address,
		},
		handlers: handlers,
	}
}

// server is implementation of core.Server for managing HTTP server lifecycle
// TODO: add mutext for running state
// TODO: add mutext for handlers
type server struct {
	running  bool
	logger   ports.Logger
	server   *http.Server
	handlers []HTTPHandler[HTTPRequest]
}

// Start starts the HTTP server
// Note: The server must be configured before starting
// If the server is already running, it will log a warning and continue
// Changes to address or handlers will take effect after restart
func (inst *server) Start() error {
	go func() {
		inst.running = true
		if err := inst.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			inst.logger.Error("server error", err)
			inst.running = false
		}
	}()

	time.Sleep(100 * time.Millisecond) // Give the server a moment to start

	if !inst.running {
		return common.ErrorServerNotRunning
	}

	return nil
}

func (inst *server) Stop() {
	if inst.server != nil || !inst.running {
		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := inst.server.Shutdown(ctx); err != nil {
			inst.logger.Error("Failed to shutdown server gracefully", err)
			// Force close if graceful shutdown fails
			if closeErr := inst.server.Close(); closeErr != nil {
				inst.logger.Error("Failed to force-close server", closeErr)
			}
		}

		inst.logger.Info("Server stopped successfully")
	}
}

func (inst *server) Restart() error {
	if inst.running {
		inst.Stop()
	}

	return inst.Start()
}

func (inst *server) IsRunning() bool {
	return inst.running
}

func (inst *server) Configure() error {
	if inst.running {
		inst.logger.Warn("Reconfiguring a running server. Changes will take effect after restart.")
	}

	// Create a new router and register handlers
	r := mux.NewRouter()

	for _, handler := range inst.handlers {
		h := handler
		methods := h.SupportedMethod()

		r.HandleFunc(h.Path(), func(w http.ResponseWriter, r *http.Request) {
			responser := &defaultResponser{}
			req, err := NewRequest(r)
			if err != nil {
				responser.Error(w, http.StatusBadRequest, "invalid request")
				return
			}

			resp, err := h.Handle(req)
			if err != nil {
				responser.Error(w, http.StatusInternalServerError, err.Error())
				return
			}

			responser.Write(w, r, resp)
		}).Methods(methods...)

		inst.logger.Info(
			"added handler:",
			zap.String("path", h.Path()),
			zap.Strings("methods", methods),
		)
	}

	inst.server.Handler = r
	return nil
}

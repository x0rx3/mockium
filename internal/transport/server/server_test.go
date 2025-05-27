package server

import (
	"mockium/internal/model"
	"mockium/internal/transport"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type MockRouter struct {
	path     string
	handlers map[model.Method]http.Handler
}

func (m *MockRouter) Path() string                            { return m.path }
func (m *MockRouter) Handler(mt model.Method) http.Handler    { return m.handlers[mt] }
func (m *MockRouter) Handlers() map[model.Method]http.Handler { return m.handlers }

func TestNewServer(t *testing.T) {
	log := zaptest.NewLogger(t)
	routes := []transport.Router{}

	srv := New(log, routes...)

	assert.NotNil(t, srv)
	assert.Equal(t, log, srv.log)
	assert.Len(t, srv.routes, 0)
}

func TestStartServer_Success(t *testing.T) {
	log := zaptest.NewLogger(t)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	routes := []transport.Router{
		&MockRouter{
			path:     "/test",
			handlers: map[model.Method]http.Handler{http.MethodGet: testHandler},
		},
	}

	srv := New(log, routes...)

	go func() {
		err := srv.Start(":0")
		assert.NoError(t, err)
	}()

	router := mux.NewRouter()
	for _, hr := range routes {
		for mth, handle := range hr.Handlers() {
			router.HandleFunc(hr.Path(), handle.ServeHTTP).Methods(string(mth))
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestStartServer_InvalidAddress(t *testing.T) {
	srv := New(zap.NewNop())

	err := srv.Start("invalid_address")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing port in address")
}

var muxNewRouter = mux.NewRouter

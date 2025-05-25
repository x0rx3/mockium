package server

import (
	"gomock/internal/transport"
	"gomock/internal/transport/method"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type MockRouter struct {
	path    string
	method  method.Method
	handler http.Handler
}

func (m *MockRouter) Path() string          { return m.path }
func (m *MockRouter) Method() method.Method { return m.method }
func (m *MockRouter) Handler() http.Handler { return m.handler }

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
			path:    "/test",
			method:  http.MethodGet,
			handler: testHandler,
		},
	}

	srv := New(log, routes...)

	go func() {
		err := srv.Start(":0")
		assert.NoError(t, err)
	}()

	router := mux.NewRouter()
	for _, hr := range routes {
		router.HandleFunc(hr.Path(), hr.Handler().ServeHTTP).Methods(string(hr.Method()))
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

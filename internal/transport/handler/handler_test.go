package handler

import (
	"encoding/json"
	"fmt"
	"gomock/internal/model"
	"gomock/internal/transport"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type MockRequestMatcher struct {
	matchFunc func(*http.Request) bool
}

func (m *MockRequestMatcher) Match(req *http.Request) bool {
	return m.matchFunc(req)
}

type MockResponseProvider struct {
	prepareFunc func(*http.Request) (*model.SetResponse, error)
}

func (m *MockResponseProvider) Prepare(req *http.Request) (*model.SetResponse, error) {
	return m.prepareFunc(req)
}

func TestNewHandler(t *testing.T) {
	log := zaptest.NewLogger(t)
	matchers := make(map[transport.RequestMatcher]transport.ResponsePreparer)

	h := New(log, matchers)

	assert.NotNil(t, h)
	assert.Equal(t, log, h.log)
	assert.Len(t, h.matchers, 0)
}

func TestServeHTTP_NotFound(t *testing.T) {
	log := zaptest.NewLogger(t)
	matchers := make(map[transport.RequestMatcher]transport.ResponsePreparer)

	h := New(log, matchers)

	req := httptest.NewRequest("GET", "/not-found", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestServeHTTP_InternalErrorOnPrepare(t *testing.T) {
	log := zaptest.NewLogger(t)

	matcher := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/error"
		},
	}

	provider := &MockResponseProvider{
		prepareFunc: func(req *http.Request) (*model.SetResponse, error) {
			return nil, fmt.Errorf("simulated error")
		},
	}

	matchers := map[transport.RequestMatcher]transport.ResponsePreparer{
		matcher: provider,
	}

	h := New(log, matchers)

	req := httptest.NewRequest("GET", "/error", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed prepare response")
}

func TestServeHTTP_JSONResponse(t *testing.T) {
	log := zaptest.NewLogger(t)
	testData := map[string]any{"key": "value"}

	matcher := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/json"
		},
	}

	provider := &MockResponseProvider{
		prepareFunc: func(req *http.Request) (*model.SetResponse, error) {
			return &model.SetResponse{
				SetBody:   testData,
				SetStatus: http.StatusCreated,
			}, nil
		},
	}

	matchers := map[transport.RequestMatcher]transport.ResponsePreparer{
		matcher: provider,
	}

	h := New(log, matchers)

	req := httptest.NewRequest("GET", "/json", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var responseBody map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Equal(t, testData, responseBody)
}

func TestServeHTTP_FileResponse(t *testing.T) {
	log := zaptest.NewLogger(t)
	testFile, err := os.Open("testdata/file.txt")
	require.NoError(t, err)

	matcher := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/file"
		},
	}

	provider := &MockResponseProvider{
		prepareFunc: func(req *http.Request) (*model.SetResponse, error) {
			return &model.SetResponse{
				SetFile:   testFile,
				SetStatus: http.StatusOK,
			}, nil
		},
	}

	matchers := map[transport.RequestMatcher]transport.ResponsePreparer{
		matcher: provider,
	}

	h := New(log, matchers)

	req := httptest.NewRequest("GET", "/file", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Disposition"), testFile.Name())
}

func TestServeHTTP_Headers(t *testing.T) {
	log := zaptest.NewLogger(t)
	testHeaders := map[string]string{"X-Test": "value"}

	matcher := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/headers"
		},
	}

	provider := &MockResponseProvider{
		prepareFunc: func(req *http.Request) (*model.SetResponse, error) {
			return &model.SetResponse{
				SetHeaders: testHeaders,
				SetStatus:  http.StatusNoContent,
			}, nil
		},
	}

	matchers := map[transport.RequestMatcher]transport.ResponsePreparer{
		matcher: provider,
	}

	h := New(log, matchers)

	req := httptest.NewRequest("GET", "/headers", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "value", rec.Header().Get("X-Test"))
}

func TestFindMatches(t *testing.T) {
	log := zaptest.NewLogger(t)

	matcher1 := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/first"
		},
	}

	matcher2 := &MockRequestMatcher{
		matchFunc: func(req *http.Request) bool {
			return req.URL.Path == "/second"
		},
	}

	provider := &MockResponseProvider{
		prepareFunc: func(req *http.Request) (*model.SetResponse, error) {
			return &model.SetResponse{}, nil
		},
	}

	matchers := map[transport.RequestMatcher]transport.ResponsePreparer{
		matcher1: provider,
		matcher2: provider,
	}

	h := New(log, matchers)

	t.Run("match first", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/first", nil)
		res := h.findMatches(req)
		assert.Equal(t, provider, res)
	})

	t.Run("match second", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/second", nil)
		res := h.findMatches(req)
		assert.Equal(t, provider, res)
	})

	t.Run("no match", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/unknown", nil)
		res := h.findMatches(req)
		assert.Nil(t, res)
	})
}

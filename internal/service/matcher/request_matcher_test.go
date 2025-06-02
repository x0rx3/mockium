package matcher

import (
	"context"
	"encoding/json"
	"mockium/internal/model"
	"mockium/internal/service/constants"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewRequestMatcher(t *testing.T) {
	tests := []struct {
		name     string
		template *model.MatchRequestTemplate
		wantErr  bool
	}{
		{
			name:     "Empty template",
			template: &model.MatchRequestTemplate{},
			wantErr:  false,
		},
		{
			name: "Template with regex",
			template: &model.MatchRequestTemplate{
				MustHeaders: map[string]any{
					"X-Request-ID": "regexp:^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid regex template",
			template: &model.MatchRequestTemplate{
				MustHeaders: map[string]any{
					"X-Request-ID": "regexp:invalid[regex",
				},
			},
			wantErr: false, // Should log warning but not fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			matcher := NewRequestMatcher(logger, tt.template)
			assert.NotNil(t, matcher)
		})
	}
}

func TestRequestMatcher_Match(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name        string
		template    *model.MatchRequestTemplate
		request     func() *http.Request
		wantMatch   bool
		wantLogMsgs []string
	}{
		{
			name: "Match simple GET request",
			template: &model.MatchRequestTemplate{
				MustPathParameters: map[string]any{
					"id": "123",
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/items/123", nil)
				req.SetPathValue("id", "123")
				return req
			},
			wantMatch: true,
		},
		{
			name: "Mismatch path parameter",
			template: &model.MatchRequestTemplate{
				MustPathParameters: map[string]any{
					"id": "456",
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/items/123", nil)
				req.SetPathValue("id", "123")
				return req
			},
			wantMatch: false,
		},
		{
			name: "Match JSON body",
			template: &model.MatchRequestTemplate{
				MustBody: map[string]any{
					"user": map[string]any{
						"name": "${regexp:^[A-Z][a-z]+$}",
						"age":  30,
					},
				},
			},
			request: func() *http.Request {

				body := `{"user": {"name": "John", "age": 30}}`
				req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
				req.Header.Add("Content-Type", constants.ContentTypeApplicationJSON)
				return req
			},
			wantMatch: true,
		},
		{
			name: "Match form data",
			template: &model.MatchRequestTemplate{
				MustBody: map[string]any{
					"username": "testuser",
					"password": "${regexp:^[a-zA-Z0-9]{8,}$}",
				},
			},
			request: func() *http.Request {
				form := url.Values{}
				form.Add("username", "testuser")
				form.Add("password", "secret123")
				req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			wantMatch: true,
		},
		{
			name: "Match query parameters",
			template: &model.MatchRequestTemplate{
				MustQueryParameters: map[string]any{
					"page":  "1",
					"limit": "${regexp:^[0-9]+$}",
				},
			},
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/items?page=1&limit=10", nil)
			},
			wantMatch: true,
		},
		{
			name: "Match headers",
			template: &model.MatchRequestTemplate{
				MustHeaders: map[string]any{
					"Content-Type": "application/json",
					"X-Request-Id": "${regexp:^[a-f0-9]{32}$}",
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("X-Request-Id", "550e8400e29b41d4a716446655440000")
				return req
			},
			wantMatch: true,
		},
		{
			name: "Body cache reuse",
			template: &model.MatchRequestTemplate{
				MustBody: map[string]any{
					"id": "${...}",
				},
			},
			request: func() *http.Request {
				body := `{"id": "123"}`
				req := httptest.NewRequest("POST", "/", strings.NewReader(body))
				req.Header.Add("Content-Type", "application/json")
				// Simulate cached body
				var mBody map[string]interface{}
				json.Unmarshal([]byte(body), &mBody)
				ctx := context.WithValue(req.Context(), ctxtBodyCacheKey{}, mBody)
				return req.WithContext(ctx)
			},
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewRequestMatcher(logger, tt.template)
			req := tt.request()
			assert.Equal(t, tt.wantMatch, matcher.Match(req))
		})
	}
}

func TestRequestMatcher_HeaderMatching(t *testing.T) {
	logger := zaptest.NewLogger(t)
	template := &model.MatchRequestTemplate{
		MustHeaders: map[string]any{
			"Authorization": "Bearer token123",
			"X-Request-Id":  "${...}",
		},
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer token123")
	req.Header.Add("X-Request-Id", "550e8400")

	matcher := NewRequestMatcher(logger, template)
	assert.True(t, matcher.Match(req))
}

func BenchmarkRequestMatcher_Match(b *testing.B) {
	logger := zaptest.NewLogger(b)
	template := &model.MatchRequestTemplate{
		MustBody: map[string]any{
			"user": map[string]any{
				"id":   "regexp:^[0-9]+$",
				"name": "regexp:^[A-Z][a-z]+$",
			},
		},
	}

	matcher := NewRequestMatcher(logger, template)
	body := `{"user": {"id": "123", "name": "John"}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		matcher.Match(req)
	}
}

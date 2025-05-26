package builder

import (
	"bytes"
	"encoding/json"
	"io"
	"mockium/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponseBuilder(t *testing.T) {
	template := model.SetResponseTemplate{}
	builder := NewResponseBuilder(template)
	assert.NotNil(t, builder)
	assert.Equal(t, template, builder.templResp)
}

func TestBuild_WithEmptyTemplate(t *testing.T) {
	template := model.SetResponseTemplate{}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := builder.Build(req)

	require.NoError(t, err)
	assert.Equal(t, 0, resp.SetStatus)
	assert.Nil(t, resp.SetHeaders)
	assert.Nil(t, resp.SetBody)
	assert.Nil(t, resp.SetFile)
}

func TestBuild_WithBodyTemplate(t *testing.T) {
	template := model.SetResponseTemplate{
		SetBody: map[string]any{
			"message": "Hello, World!",
			"nested": map[string]any{
				"key": "value",
			},
		},
		SetStatus: http.StatusOK,
		SetHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := builder.Build(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.SetStatus)
	assert.Equal(t, "application/json", resp.SetHeaders["Content-Type"])
	assert.Equal(t, "Hello, World!", resp.SetBody["message"])
	assert.Equal(t, "value", resp.SetBody["nested"].(map[string]any)["key"])
}

func TestBuild_WithFileTemplate(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("test file content")
	require.NoError(t, err)
	tmpFile.Close()

	template := model.SetResponseTemplate{
		SetFile:   tmpFile.Name(),
		SetStatus: http.StatusOK,
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := builder.Build(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.SetStatus)

	fileContent, err := io.ReadAll(resp.SetFile)
	require.NoError(t, err)
	assert.Equal(t, "test file content", string(fileContent))
}

func TestBuild_WithPlaceholders(t *testing.T) {
	tests := []struct {
		name          string
		template      model.SetResponseTemplate
		requestSetup  func(*http.Request)
		expectedValue string
		expectedKey   string
	}{
		{
			name: "Header placeholder",
			template: model.SetResponseTemplate{
				SetBody: map[string]any{
					"header": "${req.headers:X-Test-Header}",
				},
			},
			requestSetup: func(r *http.Request) {
				r.Header.Set("X-Test-Header", "test-value")
			},
			expectedValue: "test-value",
			expectedKey:   "header",
		},
		{
			name: "Query placeholder",
			template: model.SetResponseTemplate{
				SetBody: map[string]any{
					"query": "${req.query:test_param}",
				},
			},
			requestSetup: func(r *http.Request) {
				r.URL.RawQuery = "test_param=query-value"
			},
			expectedValue: "query-value",
			expectedKey:   "query",
		},
		{
			name: "Path placeholder",
			template: model.SetResponseTemplate{
				SetBody: map[string]any{
					"path": "${req.path:test_param}",
				},
			},
			requestSetup:  func(r *http.Request) {},
			expectedValue: "path-value",
			expectedKey:   "path",
		},
		{
			name: "Form placeholder",
			template: model.SetResponseTemplate{
				SetBody: map[string]any{
					"form": "${req.form:test_param}",
				},
			},
			requestSetup: func(r *http.Request) {
				r.Form = map[string][]string{"test_param": {"form-value"}}
			},
			expectedValue: "form-value",
			expectedKey:   "form",
		},
		{
			name: "Body placeholder",
			template: model.SetResponseTemplate{
				SetBody: map[string]any{
					"body": "${req.body:test_field}",
				},
			},
			requestSetup: func(r *http.Request) {
				body := map[string]string{"test_field": "body-value"}
				jsonBody, _ := json.Marshal(body)
				r.Body = io.NopCloser(bytes.NewReader(jsonBody))
			},
			expectedValue: "body-value",
			expectedKey:   "body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewResponseBuilder(tt.template)

			var req *http.Request
			if tt.name == "Path placeholder" {
				req = httptest.NewRequest("GET", "/path-value", nil)
				req = mux.SetURLVars(req, map[string]string{"test_param": "path-value"})
			} else {
				req = httptest.NewRequest("GET", "/", nil)
				if tt.requestSetup != nil {
					tt.requestSetup(req)
				}
			}

			resp, err := builder.Build(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedValue, resp.SetBody[tt.expectedKey])
		})
	}
}

func TestBuild_WithNestedPlaceholders(t *testing.T) {
	template := model.SetResponseTemplate{
		SetBody: map[string]any{
			"user": map[string]any{
				"name":     "${req.headers:X-User-Name}",
				"location": "${req.query:location}",
			},
		},
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/?location=NY", nil)
	req.Header.Set("X-User-Name", "John Doe")

	resp, err := builder.Build(req)
	require.NoError(t, err)

	user := resp.SetBody["user"].(map[string]any)
	assert.Equal(t, "John Doe", user["name"])
	assert.Equal(t, "NY", user["location"])
}

func TestBuild_WithInvalidPlaceholder(t *testing.T) {
	template := model.SetResponseTemplate{
		SetBody: map[string]any{
			"invalid": "${req.unknown:param}",
		},
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", nil)
	res, err := builder.Build(req)

	assert.NoError(t, err)
	assert.Equal(t, res.SetBody, template.SetBody)
}

func TestBuild_WithInvalidJSONBody(t *testing.T) {
	template := model.SetResponseTemplate{
		SetBody: map[string]any{
			"body": "${req.body:test_field}",
		},
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", bytes.NewReader([]byte("invalid json")))
	_, err := builder.Build(req)

	assert.Error(t, err)
}

func TestBuild_WithNonExistentFile(t *testing.T) {
	template := model.SetResponseTemplate{
		SetFile: "/nonexistent/file",
	}
	builder := NewResponseBuilder(template)

	req := httptest.NewRequest("GET", "/", nil)
	_, err := builder.Build(req)

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestValueByPlaceholders_InvalidPlaceholder(t *testing.T) {
	builder := NewResponseBuilder(model.SetResponseTemplate{})
	req := httptest.NewRequest("GET", "/", nil)

	_, err := builder.valueByPlacehoders([]string{"", "", "invalid", "param"}, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unxpected placeholder")
}

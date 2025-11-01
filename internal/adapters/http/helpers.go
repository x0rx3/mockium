package http

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

func Header(r HTTPRequest, key string) string {
	if r.Raw() == nil {
		return ""
	}
	return r.Raw().Header.Get(key)
}

func Headers(r HTTPRequest, key string) []string {
	if r.Raw() == nil {
		return nil
	}
	return r.Raw().Header.Values(key)
}

func AllHeaders(r HTTPRequest) http.Header {
	if r.Raw() == nil {
		return http.Header{}
	}
	return r.Raw().Header.Clone()
}

func JSON(r HTTPRequest, v any) error {
	body := r.Payload()
	if json.Valid(body) {
		return json.Unmarshal(body, v)
	}
	// fallback: можно пробовать читать из Raw.Body, если Payload не заполнен
	if r.Raw() != nil {
		body, _ = io.ReadAll(r.Raw().Body)
		if json.Valid(body) {
			return json.Unmarshal(body, v)
		}
	}
	return fmt.Errorf("invalid JSON payload")
}

func FormValues(r HTTPRequest) (map[string][]string, error) {
	req := r.Raw()
	if req == nil {
		return nil, fmt.Errorf("Raw() request is nil")
	}
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	return req.Form, nil
}

func Multipart(r HTTPRequest) (*multipart.Form, error) {
	req := r.Raw()
	if req == nil {
		return nil, fmt.Errorf("Raw() request is nil")
	}
	if err := req.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		return nil, err
	}
	return req.MultipartForm, nil
}

func QueryValue(r HTTPRequest, key string) string {
	if r.Raw() == nil {
		return ""
	}
	return r.Raw().URL.Query().Get(key)
}

func PathValue(r HTTPRequest, key string) string {
	if r.Raw() == nil {
		return ""
	}
	// Пример для gorilla/mux
	if vars := mux.Vars(r.Raw()); vars != nil {
		return vars[key]
	}
	return ""
}

func Cookie(r HTTPRequest, name string) (*http.Cookie, error) {
	if r.Raw() == nil {
		return nil, fmt.Errorf("Raw() request is nil")
	}
	for _, c := range r.Raw().Cookies() {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("cookie %q not found", name)
}

package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func NewRequest(r *http.Request) (HTTPRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	return &Request{
		req:     r,
		payload: body,
	}, nil
}

type Request struct {
	req     *http.Request
	payload []byte
}

// --- impl Request ---
func (inst *Request) Protocol() string                  { return "http" }
func (inst *Request) Context() context.Context          { return inst.req.Context() }
func (inst *Request) Payload() []byte                   { return inst.payload }
func (inst *Request) Metadata() map[string][]string     { return map[string][]string(inst.req.Header) }
func (inst *Request) MetadataValue(key string) []string { return inst.req.Header.Values(key) }

// --- impl HTTPRequest ---
func (inst *Request) Method() string                  { return inst.req.Method }
func (inst *Request) Path() string                    { return inst.req.URL.Path }
func (inst *Request) Raw() *http.Request              { return inst.req }
func (inst *Request) WithContext(ctx context.Context) { inst.req = inst.req.WithContext(ctx) }

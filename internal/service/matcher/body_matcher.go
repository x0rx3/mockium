package matcher

import (
	"context"
	"encoding/json"
	"io"
	"mockium/internal/service"
	"mockium/internal/service/constants"
	"net/http"

	"go.uber.org/zap"
)

// BodyMatcher is responsible for comparing HTTP request headers and body
// against expected values provided at initialization.
// It supports JSON and form-urlencoded content types.
type BodyMatcher struct {
	log          *zap.Logger      // Logger for error and debug output.
	comparer     service.Comparer // Interface for deep comparison of values.
	matchHeaders map[string]any   // Expected HTTP headers to match.
	matchBody    map[string]any   // Expected HTTP body content to match.
}

// NewBodyMatcher creates and returns a new instance of BodyMatcher.
// Parameters:
//   - log: logger used for internal logging.
//   - compare: implementation of the service.Comparer interface for matching bodies.
//   - matchHeaders: a map of expected headers (e.g., "Content-Type").
//   - matchBody: a map representing the expected structure/content of the request body.
func NewBodyMatcher(log *zap.Logger, compare service.Comparer, matchHeaders, matchBody map[string]any) *BodyMatcher {
	return &BodyMatcher{
		log:          log,
		comparer:     compare,
		matchHeaders: matchHeaders,
		matchBody:    matchBody,
	}
}

// Match checks whether the provided HTTP request satisfies the expected header and body criteria.
// Returns true if the Content-Type is correct and the request body matches the expected structure.
func (inst *BodyMatcher) Match(req *http.Request) bool {
	actualContentType := req.Header.Get("Content-Type")
	expectedContentType, ok := inst.matchHeaders["Content-Type"]

	if ok && expectedContentType != "" && actualContentType == "" {
		return false
	}

	if str, isStr := expectedContentType.(string); isStr && str != "" {
		if actualContentType == expectedContentType {
			return inst.compare(actualContentType, req)
		}
		return false
	}

	if actualContentType != "" {
		return inst.compare(actualContentType, req)
	}

	inst.log.Warn("can't parse body with empty Content-Type header")

	return false
}

// compare attempts to extract and parse the request body based on the given Content-Type,
// and compares the parsed body to the expected matchBody using the configured comparer.
// The parsed body is cached into the request's context for reuse.
func (inst *BodyMatcher) compare(headerVal string, req *http.Request) bool {
	if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).(map[string]any); cached != nil && ok {
		return inst.comparer.Compare(inst.matchBody, cached)
	}

	cached := make(map[string]any)
	switch headerVal {
	case constants.ContentTypeFormURLEncoded:
		if req.PostForm == nil {
			if err := req.ParseForm(); err != nil {
				inst.log.Error("parse form", zap.Error(err))
				return false
			}
		}

		for key, values := range req.PostForm {
			for _, value := range values {
				cached[key] = value
			}
		}

	case constants.ContentTypeApplicationJSON:
		if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).([]byte); ok {
			mBody := make(map[string]any)
			if err := json.Unmarshal(cached, &mBody); err != nil {
				inst.log.Error("parse body", zap.Error(err), zap.String("url", req.URL.Path))
				return false
			}

			return inst.comparer.Compare(inst.matchBody, mBody)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			inst.log.Warn("failed to read body", zap.String("error", err.Error()))
			return false
		}
		defer req.Body.Close()

		if err := json.Unmarshal(body, &cached); err != nil {
			inst.log.Error("parse body", zap.Error(err), zap.String("url", req.URL.Path))
			return false
		}

	default:
		inst.log.Warn("can't parse body with unexpected Content-Type header", zap.String("header", headerVal))
		return false
	}

	ctx := context.WithValue(req.Context(), ctxtBodyCacheKey{}, cached)
	*req = *req.WithContext(ctx)

	return inst.comparer.Compare(inst.matchBody, cached)
}

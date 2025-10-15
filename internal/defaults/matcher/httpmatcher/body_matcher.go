package httpmatcher

import (
	"context"
	"encoding/json"
	"io"

	"mockium/internal/defs"
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"

	"go.uber.org/zap"
)

const contentTypeHeader = "Content-Type"

type ctxtBodyCacheKey struct{}

func NewBodyMatcher(logger core.Logger, matchHeaders, matchBody map[string]any, compare core.Comparer) core.Matcher[httpgw.HTTPRequest] {
	return &bodyMatcher[httpgw.HTTPRequest]{
		logger:       logger,
		comparer:     compare,
		matchHeaders: matchHeaders,
		matchBody:    matchBody,
	}
}

type bodyMatcher[T httpgw.HTTPRequest] struct {
	logger       core.Logger
	comparer     core.Comparer
	matchHeaders map[string]any
	matchBody    map[string]any
}

func (inst *bodyMatcher[T]) Match(req T) bool {
	actualContentType := req.Raw().Header.Get(contentTypeHeader)
	expectedContentType, ok := inst.matchHeaders[contentTypeHeader]

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

	inst.logger.Warn("can't parse body with empty Content-Type header")

	return false
}

func (inst *bodyMatcher[T]) compare(headerVal string, req T) bool {
	if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).(map[string]any); cached != nil && ok {
		return inst.comparer.Compare(inst.matchBody, cached)
	}

	cached := make(map[string]any)
	switch headerVal {
	case defs.ContentTypeFormURLEncoded:
		if req.Raw().PostForm == nil {
			if err := req.Raw().ParseForm(); err != nil {
				inst.logger.Error("parse form", err)
				return false
			}
		}

		for key, values := range req.Raw().PostForm {
			for _, value := range values {
				cached[key] = value
			}
		}

	case defs.ContentTypeApplicationJSON:
		if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).([]byte); ok {
			mBody := make(map[string]any)
			if err := json.Unmarshal(cached, &mBody); err != nil {
				inst.logger.Error("parse body"+req.Raw().URL.Path, err)
				return false
			}

			return inst.comparer.Compare(inst.matchBody, mBody)
		}

		body, err := io.ReadAll(req.Raw().Body)
		if err != nil {
			inst.logger.Warn("failed to read body", zap.String("error", err.Error()))
			return false
		}
		defer req.Raw().Body.Close()

		if err := json.Unmarshal(body, &cached); err != nil {
			inst.logger.Error("parse body"+req.Raw().URL.Path, err)
			return false
		}

	default:
		inst.logger.Warn("can't parse body with unexpected Content-Type header", zap.String("header", headerVal))
		return false
	}

	req.WithContext(context.WithValue(req.Context(), ctxtBodyCacheKey{}, cached))

	return inst.comparer.Compare(inst.matchBody, cached)
}

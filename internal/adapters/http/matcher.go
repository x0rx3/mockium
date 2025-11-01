package http

import (
	"context"
	"encoding/json"
	"io"
	"mockium/internal/adapters/common"
	"mockium/pkg/model"
	"mockium/pkg/ports"

	"go.uber.org/zap"
)

const cookieHeaderKey = "Cookie"
const contentTypeHeader = "Content-Type"

type ctxtBodyCacheKey struct{}

func NewRequestMathcer(logger ports.Logger, handle model.Handle, comparer ports.Comparer) ports.Matcher[HTTPRequest] {
	mathcers := make([]ports.Matcher[HTTPRequest], 0)

	otherHeaders := make(map[string]any, 0)
	if handle.MatchRequest.MustHeaders != nil {
		mathcers = append(
			mathcers,
			NewHeadersMatcher(handle.MatchRequest.MustHeaders, comparer),
		)
		otherHeaders = handle.MatchRequest.MustHeaders.Other
	}

	if handle.MatchRequest.MustBody != nil {
		mathcers = append(
			mathcers,
			NewBodyMatcher(
				logger,
				otherHeaders,
				handle.MatchRequest.MustBody,
				comparer,
			),
		)
	}

	if handle.MatchRequest.MustPathParameters != nil {
		mathcers = append(
			mathcers,
			NewPathMatcher(
				handle.MatchRequest.MustPathParameters,
				comparer,
			),
		)
	}

	if handle.MatchRequest.MustQueryParameters != nil {
		mathcers = append(
			mathcers,
			NewQueryMatcher(
				handle.MatchRequest.MustQueryParameters,
				comparer,
			),
		)
	}

	return &requestMatcher{
		mathcers: mathcers,
	}
}

type requestMatcher struct {
	mathcers []ports.Matcher[HTTPRequest]
}

func (inst *requestMatcher) Match(req HTTPRequest) bool {
	for _, mathcer := range inst.mathcers {
		if !mathcer.Match(req) {
			return false
		}
	}

	return true
}

func NewHeadersMatcher(matchHeaders *model.MustHeader, comparer ports.Comparer) ports.Matcher[HTTPRequest] {
	return &headersMatcher{
		matchHeaders: *matchHeaders,
		comparer:     comparer,
	}
}

type headersMatcher struct {
	matchHeaders model.MustHeader // Expected headers to be matched.
	comparer     ports.Comparer   // Comparer used to check header values.
}

func (inst *headersMatcher) Match(req HTTPRequest) bool {
	if len(inst.matchHeaders.Other) != 0 {
		for key, tValue := range inst.matchHeaders.Other {
			actual := req.Raw().Header.Get(key)
			if actual == "" || !inst.comparer.Compare(tValue, actual) {
				return false
			}
		}
	}

	if len(inst.matchHeaders.Cookie) != 0 {
		if !inst.compareCookie(inst.matchHeaders.Cookie, req) {
			return false
		}
	}

	return true
}

// compareCookie
func (inst *headersMatcher) compareCookie(expCookies []model.Cookie, req HTTPRequest) bool {
	actualCookies := req.Raw().Cookies()
	if len(actualCookies) == 0 && len(expCookies) == 0 {
		return true
	}

	if len(actualCookies) == 0 && len(expCookies) > 0 {
		return false
	}

	if len(actualCookies) != 0 && len(expCookies) == 0 {
		return true
	}

	for _, expCookie := range expCookies {
		if expCookie.Name == "" {
			for _, actCookie := range req.Raw().Cookies() {
				if inst.comparer.Compare(expCookie.Value, actCookie.Value) {
					return true
				}
			}
			return false
		}

		actualCookie, _ := req.Raw().Cookie(expCookie.Name)

		if !inst.comparer.Compare(expCookie.Value, actualCookie.Value) {
			return false
		}
	}

	return true
}

func NewQueryMatcher(matchQuery map[string]any, comparer ports.Comparer) ports.Matcher[HTTPRequest] {
	return &queryMatcher{
		matchQuery: matchQuery,
		comparer:   comparer,
	}
}

type queryMatcher struct {
	matchQuery map[string]any
	comparer   ports.Comparer
}

func (inst *queryMatcher) Match(req HTTPRequest) bool {
	query := req.Raw().URL.Query()

	for expKey, expVal := range inst.matchQuery {
		actVals, ok := query[expKey]
		if !ok {
			return false // ключа нет
		}

		match := false
		for _, v := range actVals {
			if v == expVal {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}

func NewPathMatcher(matchPath map[string]any, comparer ports.Comparer) ports.Matcher[HTTPRequest] {
	return &pathMatcher{
		matchPath: matchPath,
		comparer:  comparer,
	}
}

type pathMatcher struct {
	matchPath map[string]any
	comparer  ports.Comparer
}

func (inst *pathMatcher) Match(req HTTPRequest) bool {
	for key, tValue := range inst.matchPath {
		actual := req.Raw().PathValue(key)
		if actual == "" || !inst.comparer.Compare(tValue, actual) {
			return false
		}
	}
	return true
}

func NewBodyMatcher(logger ports.Logger, matchHeaders, matchBody map[string]any, compare ports.Comparer) ports.Matcher[HTTPRequest] {
	return &bodyMatcher{
		logger:       logger,
		comparer:     compare,
		matchHeaders: matchHeaders,
		matchBody:    matchBody,
	}
}

type bodyMatcher struct {
	logger       ports.Logger
	comparer     ports.Comparer
	matchHeaders map[string]any
	matchBody    map[string]any
}

func (inst *bodyMatcher) Match(req HTTPRequest) bool {
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

func (inst *bodyMatcher) compare(headerVal string, req HTTPRequest) bool {
	if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).(map[string]any); cached != nil && ok {
		return inst.comparer.Compare(inst.matchBody, cached)
	}

	cached := make(map[string]any)
	switch headerVal {
	case common.ContentTypeFormURLEncoded:
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

	case common.ContentTypeApplicationJSON:
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

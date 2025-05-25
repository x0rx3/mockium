package matcher

import (
	"context"
	"encoding/json"
	"gomock/internal/model"
	"gomock/internal/service"
	"io"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

type ctxtBodyCacheKey struct{}

// RequestMatcher matches HTTP requests against predefined templates.
// It supports matching of form parameters, path parameters, query parameters,
// headers, and body with regex support.
type RequestMatcher struct {
	log          *zap.Logger
	matchRequest model.MatchRequest
	matchFuncs   []func(req *http.Request) bool
}

// NewRequestMatcher creates a new RequestMatcher instance with the specified logger and template.
// It precompiles all matching rules from the template including regex patterns.
// The matcher is ready to use immediately after creation.
func NewRequestMatcher(log *zap.Logger, templateRequest *model.MatchRequestTemplate) *RequestMatcher {
	requestMatcher := &RequestMatcher{
		log: log,
	}

	matchRequest := model.MatchRequest{
		MustFormParameters:  requestMatcher.precompile(service.Form, templateRequest.MustFormParameters),
		MustPathParameters:  requestMatcher.precompile(service.Path, templateRequest.MustPathParameters),
		MustQueryParameters: requestMatcher.precompile(service.Query, templateRequest.MustQueryParameters),
		MustBody:            requestMatcher.precompile(service.Body, templateRequest.MustBody),
		MustHeaders:         requestMatcher.precompile(service.Headers, templateRequest.MustHeaders),
	}

	requestMatcher.matchRequest = matchRequest

	return requestMatcher
}

// Match checks if the provided HTTP request matches all configured criteria.
// Returns true if all match conditions are satisfied, false otherwise.
// The matching is performed in the order: form, path, query, body, headers.
func (inst *RequestMatcher) Match(req *http.Request) bool {
	for _, match := range inst.matchFuncs {
		if !match(req) {
			return false
		}
	}
	return true
}

// precompile processes matching rules for a specific parameter type.
// It adds the corresponding match function to the matcher and compiles any regex patterns.
// Returns the processed matching rules with compiled regex patterns.
func (inst *RequestMatcher) precompile(param service.Parameter, source map[string]any) map[string]any {
	if len(source) == 0 {
		return nil
	}

	inst.addMatchFunc(param)
	return inst.precompileRegexp(source)
}

// precompileRegexp compiles all regex patterns in the provided source map.
// Returns a new map with string values replaced by compiled regex patterns where applicable.
func (inst *RequestMatcher) precompileRegexp(source map[string]any) map[string]any {
	result := make(map[string]any)
	for key, value := range source {
		if str, ok := value.(string); ok {
			if service.RegexpResponseValuePlaceholder.MatchString(str) {
				placeholders := service.RegexpResponseValuePlaceholder.FindStringSubmatch(str)
				switch placeholders[2] {
				case service.RegexpValuePlaceholder:
					if re, err := regexp.Compile(placeholders[3]); err == nil {
						result[key] = re
					} else {
						inst.log.Warn("failed to compile regexp", zap.String("regexp", str))
						result[key] = str
					}
				default:
					result[key] = str
				}
			}
			continue
		}
		result[key] = value
	}
	return result
}

// addMatchFunc registers the appropriate matching function based on parameter type.
func (inst *RequestMatcher) addMatchFunc(param service.Parameter) {
	switch param {
	case service.Form:
		inst.matchFuncs = append(inst.matchFuncs, inst.matchForm)
	case service.Path:
		inst.matchFuncs = append(inst.matchFuncs, inst.matchPath)
	case service.Query:
		inst.matchFuncs = append(inst.matchFuncs, inst.matchQuery)
	case service.Body:
		inst.matchFuncs = append(inst.matchFuncs, inst.matchBody)
	case service.Headers:
		inst.matchFuncs = append(inst.matchFuncs, inst.matchHeader)
	}
}

// matchForm checks if the request's form parameters match the template.
func (inst *RequestMatcher) matchForm(req *http.Request) bool {
	if req.PostForm == nil {
		if err := req.ParseForm(); err != nil {
			inst.log.Debug("failed to parse form", zap.Error(err))
			return false
		}
	}

	for key, tValue := range inst.matchRequest.MustFormParameters {
		actual := req.PostForm.Get(key)
		if actual != "" || !inst.compare(tValue, actual) {
			return false
		}
	}
	return true
}

// matchPath checks if the request's path parameters match the template.
func (inst *RequestMatcher) matchPath(req *http.Request) bool {
	for key, tValue := range inst.matchRequest.MustPathParameters {
		actual := req.PathValue(key)
		if actual == "" || !inst.compare(tValue, actual) {
			return false
		}
	}
	return true
}

// matchQuery checks if the request's query parameters match the template.
func (inst *RequestMatcher) matchQuery(req *http.Request) bool {
	for key, tValue := range inst.matchRequest.MustQueryParameters {
		actual := req.URL.Query().Get(key)
		if actual == "" || !inst.compare(tValue, actual) {
			return false
		}
	}
	return true
}

// matchHeader checks if the request's headers match the template.
func (inst *RequestMatcher) matchHeader(req *http.Request) bool {
	return false
}

// matchBody checks if the request's body matches the template.
// Caches the parsed body in the request context for subsequent matches.
func (inst *RequestMatcher) matchBody(req *http.Request) bool {
	if cached, ok := req.Context().Value(ctxtBodyCacheKey{}).(map[string]interface{}); ok {
		return inst.compareBody(cached)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		inst.log.Warn("failed to read body", zap.String("error", err.Error()))
		return false
	}
	defer req.Body.Close()

	mBody := make(map[string]any)
	if err := json.Unmarshal(body, &mBody); err != nil {
		return false
	}

	ctx := context.WithValue(req.Context(), ctxtBodyCacheKey{}, mBody)
	*req = *req.WithContext(ctx)

	return inst.compareBody(mBody)
}

// compare checks if the actual value matches the expected pattern.
// Supports direct comparison, regex matching (for strings), and deep comparison of maps and slices.
func (inst *RequestMatcher) compare(expected, actual any) bool {
	switch exp := expected.(type) {
	case *regexp.Regexp:
		if str, ok := actual.(string); ok {
			return exp.MatchString(str)
		}
		return false
	case string:
		if exp == service.AnyValuePlaceholder {
			return true
		}
		return exp == actual
	case []any:
		if aSlice, ok := actual.([]any); ok {
			return inst.compareSlices(exp, aSlice)
		}
		return false
	case map[string]any:
		if aMap, ok := actual.(map[string]any); ok {
			return inst.compareMaps(exp, aMap)
		}
		return false
	}
	return expected == actual
}

// compareBody checks if the actual body matches all required body fields from the template.
func (inst *RequestMatcher) compareBody(actual map[string]any) bool {
	for tKey, tValue := range inst.matchRequest.MustBody {
		if v, ok := actual[tKey]; !ok || !inst.compare(tValue, v) {
			return false
		}
	}
	return true
}

// compareSlices checks if two slices are deeply equal according to the matcher's comparison rules.
func (inst *RequestMatcher) compareSlices(expected, actual []interface{}) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := range expected {
		if !inst.compare(expected[i], actual[i]) {
			return false
		}
	}
	return true
}

// compareMaps checks if two maps are deeply equal according to the matcher's comparison rules.
func (inst *RequestMatcher) compareMaps(expected, actual map[string]interface{}) bool {
	for key, expVal := range expected {
		actVal, exists := actual[key]
		if !exists {
			return false
		}
		if !inst.compare(expVal, actVal) {
			return false
		}
	}
	return true
}

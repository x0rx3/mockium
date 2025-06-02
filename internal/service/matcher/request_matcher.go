package matcher

import (
	"mockium/internal/model"
	"mockium/internal/service/comparer"
	"mockium/internal/service/constants"
	"mockium/internal/transport"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

// ctxtBodyCacheKey is an unexported struct used as a context key
// for caching parsed request body data to avoid redundant parsing.
type ctxtBodyCacheKey struct{}

// RequestMatcher aggregates multiple request matchers (e.g., body, headers, path, query)
// and evaluates an HTTP request against all of them.
type RequestMatcher struct {
	log               *zap.Logger                // Logger for diagnostic messages.
	parameterMatchers []transport.RequestMatcher // Set of matchers to evaluate against the request.
}

// NewRequestMatcher creates a new RequestMatcher based on a template request specification.
// It compiles matchers for path parameters, headers, body, and query parameters as needed.
//
// Parameters:
//   - log: a structured logger used for diagnostics.
//   - templateRequest: the specification that defines required parameters for a match.
//
// Returns a pointer to a configured RequestMatcher instance.
func NewRequestMatcher(log *zap.Logger, templateRequest *model.MatchRequestTemplate) *RequestMatcher {
	requestMatcher := &RequestMatcher{
		log: log,
	}

	comparer := comparer.New()

	parameterMatchers := make([]transport.RequestMatcher, 0)
	if len(templateRequest.MustPathParameters) > 0 {
		parameterMatchers = append(parameterMatchers, NewPathMatcher(requestMatcher.precompileRegexp(templateRequest.MustPathParameters), comparer))
	}

	matchHeaders := make(map[string]any)
	if len(templateRequest.MustHeaders) > 0 {
		matchHeaders = requestMatcher.precompileRegexp(templateRequest.MustHeaders)
		parameterMatchers = append(parameterMatchers, NewHeadersMatcher(matchHeaders, comparer))
	}

	if len(templateRequest.MustBody) > 0 {
		parameterMatchers = append(parameterMatchers, NewBodyMatcher(log, comparer, matchHeaders, requestMatcher.precompileRegexp(templateRequest.MustBody)))
	}

	if len(templateRequest.MustQueryParameters) > 0 {
		parameterMatchers = append(parameterMatchers, NewQueryMatcher(requestMatcher.precompileRegexp(templateRequest.MustQueryParameters), comparer))
	}

	if len(templateRequest.MustPathParameters) > 0 {
		parameterMatchers = append(parameterMatchers, NewPathMatcher(requestMatcher.precompileRegexp(templateRequest.MustPathParameters), comparer))
	}

	requestMatcher.parameterMatchers = parameterMatchers

	return requestMatcher
}

// Match runs all internal matchers against the given HTTP request.
// The request must pass all matchers to be considered a match.
//
// Returns true if all parameter matchers validate successfully; false otherwise.
func (inst *RequestMatcher) Match(req *http.Request) bool {
	for _, match := range inst.parameterMatchers {
		if !match.Match(req) {
			return false
		}
	}
	return true
}

// precompileRegexp recursively processes a map of values that may include
// regular expression placeholders. If a placeholder is detected, it attempts
// to compile it into a *regexp.Regexp object.
//
// Returns a new map with compiled regexes where applicable.
// If compilation fails, the original string is retained and a warning is logged.
func (inst *RequestMatcher) precompileRegexp(source map[string]any) map[string]any {
	if len(source) == 0 {
		return nil
	}

	result := make(map[string]any, len(source))
	for key, value := range source {
		switch v := value.(type) {
		case string:
			if constants.RegexpRequestValuePlaceholder.MatchString(v) {
				placeholders := constants.RegexpRequestValuePlaceholder.FindStringSubmatch(v)
				if placeholders[1] == constants.RegexpValuePlaceholder {
					if re, err := regexp.Compile(placeholders[2]); err == nil {
						result[key] = re
						continue
					} else {
						inst.log.Warn("failed to compile regexp", zap.String("regexp", v))
						result[key] = v
						continue
					}
				}
			}
			result[key] = v
		case map[string]any:
			result[key] = inst.precompileRegexp(v)
		default:
			result[key] = value
		}
	}
	return result
}

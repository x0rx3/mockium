package matcher

import (
	"mockium/internal/service"
	"net/http"
)

// HeadersMatcher is responsible for validating whether an HTTP request's headers
// match a predefined set of expected values.
type HeadersMatcher struct {
	matchHeaders map[string]any   // Expected headers to be matched.
	comparer     service.Comparer // Comparer used to check header values.
}

// NewHeadersMatcher returns a new instance of HeadersMatcher.
//
// Parameters:
//   - matchHeaders: a map of expected header key-value pairs.
//   - comparer: an implementation of service.Comparer to compare actual and expected header values.
func NewHeadersMatcher(matchHeaders map[string]any, comparer service.Comparer) *HeadersMatcher {
	return &HeadersMatcher{
		matchHeaders: matchHeaders,
		comparer:     comparer,
	}
}

// Match evaluates whether all expected headers are present in the provided HTTP request,
// and whether their values match according to the configured comparer.
//
// Returns true if all expected headers are found and match; otherwise, returns false.
func (inst *HeadersMatcher) Match(req *http.Request) bool {
	for key, tValue := range inst.matchHeaders {
		actual := req.Header.Get(key)
		if actual == "" || !inst.comparer.Compare(tValue, actual) {
			return false
		}
	}
	return true
}

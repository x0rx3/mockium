package matcher

import (
	"mockium/internal/service"
	"net/http"
)

// QueryMatcher checks whether the query parameters of an HTTP request
// match a predefined set of expected values.
type QueryMatcher struct {
	matchQuery map[string]any   // Expected query parameters to match.
	comparer   service.Comparer // Comparer used to evaluate parameter values.
}

// NewQueryMatcher creates and returns a new instance of QueryMatcher.
//
// Parameters:
//   - matchQuery: a map of expected query parameters and their values.
//   - comparer: an implementation of service.Comparer used to compare actual vs. expected values.
func NewQueryMatcher(matchQuery map[string]any, comparer service.Comparer) *QueryMatcher {
	return &QueryMatcher{
		matchQuery: matchQuery,
		comparer:   comparer,
	}
}

// Match determines whether all expected query parameters are present in the HTTP request
// and whether their values match the expected values using the configured comparer.
//
// Returns true if all query parameters match; otherwise, returns false.
func (inst *QueryMatcher) Match(req *http.Request) bool {
	for key, tValue := range inst.matchQuery {
		actual := req.URL.Query().Get(key)
		if actual == "" || !inst.comparer.Compare(tValue, actual) {
			return false
		}
	}
	return true
}

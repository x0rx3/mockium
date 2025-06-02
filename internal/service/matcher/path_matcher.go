package matcher

import (
	"mockium/internal/service"
	"net/http"
)

// PathMatcher is responsible for checking whether path parameters from an HTTP request
// match the expected values defined in the matcher.
type PathMatcher struct {
	matchPath map[string]any   // Expected path parameters to match.
	comparer  service.Comparer // Comparer used to compare actual and expected values.
}

// NewPathMatcher creates and returns a new instance of PathMatcher.
//
// Parameters:
//   - matchPath: a map of expected path parameters (e.g., route variables).
//   - comparer: an implementation of service.Comparer to evaluate matches.
func NewPathMatcher(matchPath map[string]any, comparer service.Comparer) *PathMatcher {
	return &PathMatcher{
		matchPath: matchPath,
		comparer:  comparer,
	}
}

// Match checks whether all expected path parameters exist in the request and match the
// expected values using the configured comparer.
//
// Returns true if all path values match; otherwise, returns false.
func (inst *PathMatcher) Match(req *http.Request) bool {
	for key, tValue := range inst.matchPath {
		actual := req.PathValue(key)
		if actual == "" || !inst.comparer.Compare(tValue, actual) {
			return false
		}
	}
	return true
}

package httpmatcher

import (
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
)

func NewPathMatcher(matchPath map[string]any, comparer core.Comparer) core.Matcher[httpgw.HTTPRequest] {
	return &pathMatcher[httpgw.HTTPRequest]{
		matchPath: matchPath,
		comparer:  comparer,
	}
}

type pathMatcher[T httpgw.HTTPRequest] struct {
	matchPath map[string]any
	comparer  core.Comparer
}

func (inst *pathMatcher[T]) Match(req T) bool {
	for key, tValue := range inst.matchPath {
		actual := req.Raw().PathValue(key)
		if actual == "" || !inst.comparer.Compare(tValue, actual) {
			return false
		}
	}
	return true
}

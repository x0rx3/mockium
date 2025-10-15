package httpmatcher

import (
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
)

func NewQueryMatcher(matchQuery map[string]any, comparer core.Comparer) core.Matcher[httpgw.HTTPRequest] {
	return &queryMatcher[httpgw.HTTPRequest]{
		matchQuery: matchQuery,
		comparer:   comparer,
	}
}

type queryMatcher[T httpgw.HTTPRequest] struct {
	matchQuery map[string]any
	comparer   core.Comparer
}

func (inst *queryMatcher[T]) Match(req T) bool {
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

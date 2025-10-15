package httpmatcher

// import (
// 	"mockium/internal/gateway/httpgw"
// 	"mockium/internal/model"
// 	"mockium/internal/service/compare"
// 	"mockium/internal/service/match"

// 	"go.uber.org/zap"
// )

// type ctxtBodyCacheKey struct{}

// func NewMatcher(log *zap.Logger, template *model.MatchRequest) match.Matcher[httpgw.HTTPRequest] {
// 	inst := &requestMatcher[httpgw.HTTPRequest]{
// 		log: log,
// 	}

// 	comparer := compare.New()

// 	if len(template.MustQueryParameters) > 0 {
// 		inst.matchers = append(inst.matchers, NewQueryMatcher(template.MustQueryParameters, comparer))
// 	}

// 	if len(template.MustPathParameters) > 0 {
// 		inst.matchers = append(inst.matchers, NewPathMatcher(template.MustPathParameters, comparer))
// 	}

// 	if template.MustHeaders.Cookie != nil || len(template.MustHeaders.Other) > 0 {
// 		inst.matchers = append(inst.matchers, NewHeadersMatcher(log, template.MustHeaders, comparer))

// 	}

// 	if len(template.MustBody) > 0 {
// 		inst.matchers = append(inst.matchers, NewBodyMatcher(log, template.MustHeaders.Other, template.MustBody, comparer))
// 	}

// 	return inst
// }

// type requestMatcher[T httpgw.HTTPRequest] struct {
// 	log      *zap.Logger
// 	matchers []match.Matcher[T]
// }

// func (inst *requestMatcher[T]) Match(req T) bool {
// 	if len(inst.matchers) == 0 {
// 		return true
// 	}

// 	for _, m := range inst.matchers {
// 		if !m.Match(req) {
// 			return false
// 		}
// 	}

// 	return true
// }

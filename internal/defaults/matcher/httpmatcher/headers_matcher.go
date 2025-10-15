package httpmatcher

import (
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

const cookieHeaderKey = "Cookie"

func NewHeadersMatcher(matchHeaders model.MustHeader, comparer core.Comparer) core.Matcher[httpgw.HTTPRequest] {
	return &headersMatcher[httpgw.HTTPRequest]{
		matchHeaders: matchHeaders,
		comparer:     comparer,
	}
}

type headersMatcher[T httpgw.HTTPRequest] struct {
	matchHeaders model.MustHeader // Expected headers to be matched.
	comparer     core.Comparer    // Comparer used to check header values.
}

func (inst *headersMatcher[T]) Match(req T) bool {
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
func (inst *headersMatcher[T]) compareCookie(expCookies []model.Cookie, req T) bool {
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

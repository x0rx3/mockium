package httpbld

import (
	"mockium/internal/defaults/comparer"
	"mockium/internal/defaults/matcher/httpmatcher"
	"mockium/internal/defaults/service/httpsrv"
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

var BuildHandler = func(logger core.Logger, path string, handles []model.Handle) httpgw.HTTPHandler[httpgw.HTTPRequest] {
	mthSrv := make(map[string]core.Service[httpgw.HTTPRequest], 0)
	cmpr := comparer.New()
	for _, handle := range handles {
		mthSrv[handle.MatchRequest.MustMethod] =
			httpsrv.NewService(
				core.ResMtchMap[httpgw.HTTPRequest]{
					NewResponseBuilder(&handle.SetResponse): []core.Matcher[httpgw.HTTPRequest]{
						httpmatcher.NewBodyMatcher(
							logger,
							handle.MatchRequest.MustHeaders.Other,
							handle.MatchRequest.MustBody,
							cmpr,
						),
						httpmatcher.NewHeadersMatcher(
							handle.MatchRequest.MustHeaders,
							cmpr,
						),
						httpmatcher.NewPathMatcher(
							handle.MatchRequest.MustPathParameters,
							cmpr,
						),
						httpmatcher.NewQueryMatcher(
							handle.MatchRequest.MustQueryParameters,
							cmpr,
						),
					},
				},
			)
	}

	return httpgw.NewHandler(path, mthSrv)
}

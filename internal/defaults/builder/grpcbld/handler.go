package grpcbld

import (
	"mockium/internal/defaults/matcher/grpcmatcher"
	"mockium/internal/defaults/service/grpcsrv"
	"mockium/internal/gateway/grpcgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

var BuildHandler = func(logger core.Logger, service string, handles []model.Handle) grpcgw.GRPCHandler[grpcgw.GRPCRequest] {
	unarySrv := make(map[string]core.Service[grpcgw.GRPCRequest], 0)
	steamSrv := make(map[string]core.Service[grpcgw.GRPCRequest], 0)

	for _, handle := range handles {
		switch handle.MatchRequest.MustType {
		case "unary":
			unarySrv[handle.MatchRequest.MustMethod] =
				grpcsrv.NewService(
					logger,
					core.ResMtchMap[grpcgw.GRPCRequest]{
						NewResponseBuilder(&handle.SetResponse): []core.Matcher[grpcgw.GRPCRequest]{
							grpcmatcher.NewRequestMatcher(),
						},
					},
				)
		case "stream":

		}
	}

	return grpcgw.NewHandler(
		service,
		unarySrv,
		steamSrv,
	)
}

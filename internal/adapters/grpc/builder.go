package grpc

import (
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

var BuildHandler = func(logger ports.Logger, service string, handles []model.Handle) GRPCHandler[GRPCRequest] {
	unarySrv := make(map[string]ports.Service[GRPCRequest], 0)
	steamSrv := make(map[string]ports.Service[GRPCRequest], 0)

	for _, handle := range handles {
		switch handle.MatchRequest.MustType {
		case "unary":
			unarySrv[handle.MatchRequest.MustMethod] =
				NewService(
					logger,
					ports.ResMtchMap[GRPCRequest]{
						NewResponseBuilder(&handle.SetResponse): NewRequestMatcher(),
					},
				)
		case "stream":

		}
	}

	return NewHandler(
		service,
		unarySrv,
		steamSrv,
	)
}

func NewResponseBuilder(setResponse *model.SetResponse) ports.ResponseBuilder[GRPCRequest] {
	return &responseBuilder{}
}

type responseBuilder struct {
	// no fields needed
}

func (inst *responseBuilder) Build(req GRPCRequest) (*model.SetResponse, error) {
	// Реализация построения ответа на основе запроса
	return &model.SetResponse{
		SetStatus: 200,
		SetBody:   map[string]any{"message": "Hello, gRPC!"},
	}, nil
}

func NewServerBuilder(
	config *ports.ServerConfig,
) ports.ServerBuilder {
	return &serverBuilder{
		config: config,
	}
}

type serverBuilder struct {
	config *ports.ServerConfig
}

func (inst *serverBuilder) Build(logger ports.Logger, templs []model.Template[[]model.Handle]) ports.Server {
	return NewServer(logger, inst.config.IP+":"+inst.config.Port, nil)
}

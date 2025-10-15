package grpcbld

import (
	"mockium/internal/gateway/grpcgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func NewResponseBuilder(setResponse *model.SetResponse) core.ResponseBuilder[grpcgw.GRPCRequest] {
	return &responseBuilder{}
}

type responseBuilder struct {
	// no fields needed
}

func (inst *responseBuilder) Build(req grpcgw.GRPCRequest) (*model.SetResponse, error) {
	// Реализация построения ответа на основе запроса
	return &model.SetResponse{
		SetStatus: 200,
		SetBody:   map[string]any{"message": "Hello, gRPC!"},
	}, nil
}

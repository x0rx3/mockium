package grpcgw

import (
	"context"
	"mockium/internal/gateway/defaultgw"
	"mockium/pkg/core"
)

func NewRequest(
	ctx context.Context,
	method string,
	reqType string,
	payload []byte,
	metadata map[string][]string,
) GRPCRequest {
	return &request{
		Request: defaultgw.NewRequest(ctx, "grpc", payload, metadata),
		method:  method,
		reqType: reqType,
	}
}

type request struct {
	core.Request
	method  string
	reqType string
}

func (inst *request) Method() string { return inst.method }

func (inst *request) Type() string { return inst.reqType }

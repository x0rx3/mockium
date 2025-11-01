package grpc

import (
	"context"
	"mockium/internal/adapters/common"
	"mockium/pkg/ports"
)

func NewRequest(
	ctx context.Context,
	method string,
	reqType string,
	payload []byte,
	metadata map[string][]string,
) GRPCRequest {
	return &request{
		Request: common.NewRequest(ctx, "grpc", payload, metadata),
		method:  method,
		reqType: reqType,
	}
}

type request struct {
	ports.Request
	method  string
	reqType string
}

func (inst *request) Method() string { return inst.method }

func (inst *request) Type() string { return inst.reqType }

package grpcgw

import (
	"mockium/pkg/core"
)

type Method struct {
	name string
	typ  string
}

type GRPCRequest interface {
	core.Request
	Method() string
	Type() string
}

type GRPCHandler[T GRPCRequest] interface {
	core.Handler[T]
	Service() string
	StreamMethods() []string
	UnaryMethods() []string
}

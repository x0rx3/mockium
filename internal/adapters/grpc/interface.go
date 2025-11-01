package grpc

import (
	"mockium/pkg/ports"
)

type Method struct {
	name string
	typ  string
}

type GRPCRequest interface {
	ports.Request
	Method() string
	Type() string
}

type GRPCHandler[T GRPCRequest] interface {
	ports.Handler[T]
	Service() string
	StreamMethods() []string
	UnaryMethods() []string
}

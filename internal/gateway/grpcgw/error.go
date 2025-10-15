package grpcgw

import "errors"

var (
	errInvalidMethod = errors.New("invalid gRPC method")
	errNoMatch       = errors.New("no found")
)

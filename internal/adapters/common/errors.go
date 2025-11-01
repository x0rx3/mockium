package common

import "errors"

var (
	ErrorUnknownHandleType   = errors.New("unknown error type")
	ErrorUnxpectedHandleType = errors.New("unexpected template type")
	ErrorSetBodyWithSetFile  = errors.New("cannot use parameter 'SetBody' with 'SetFile'")
	ErrorUnxpectedHTTPMethod = errors.New("unxpected http method")
	ErrorServerNotRunning    = errors.New("server is not running")
)

var (
	ErrInvalidMethod = errors.New("invalid gRPC method")
	ErrNoMatch       = errors.New("no found")
)

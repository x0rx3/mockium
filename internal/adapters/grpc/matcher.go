package grpc

import (
	"mockium/pkg/ports"
)

func NewRequestMatcher() ports.Matcher[GRPCRequest] {
	return &requestMatcher{}
}

type requestMatcher struct {
	// no fields needed
}

func (inst *requestMatcher) Match(req GRPCRequest) bool {
	// Реализация логики сопоставления запроса
	return true
}

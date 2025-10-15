package grpcmatcher

import (
	"mockium/internal/gateway/grpcgw"
	"mockium/pkg/core"
)

func NewRequestMatcher() core.Matcher[grpcgw.GRPCRequest] {
	return &requestMatcher{}
}

type requestMatcher struct {
	// no fields needed
}

func (inst *requestMatcher) Match(req grpcgw.GRPCRequest) bool {
	// Реализация логики сопоставления запроса
	return true
}

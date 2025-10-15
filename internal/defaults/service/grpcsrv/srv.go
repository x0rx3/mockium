package grpcsrv

import (
	"fmt"
	"mockium/internal/gateway/grpcgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func NewService(
	logger core.Logger,
	m core.ResMtchMap[grpcgw.GRPCRequest],
) core.Service[grpcgw.GRPCRequest] {
	return &service{
		logger,
		m,
	}
}

type service struct {
	logger core.Logger
	m      core.ResMtchMap[grpcgw.GRPCRequest]
}

func (inst *service) Handle(req grpcgw.GRPCRequest) (*model.SetResponse, error) {
	resProvider := inst.findMatches(req)
	if resProvider == nil {
		return nil, fmt.Errorf("not found")
	}

	response, _ := resProvider.Build(req)
	if response == nil {
		return nil, fmt.Errorf("nil response after prepare")
	}

	return response, nil

}

func (inst *service) findMatches(req grpcgw.GRPCRequest) core.ResponseBuilder[grpcgw.GRPCRequest] {
	for resBld, reqMtchs := range inst.m {
		find := true
		for _, reqMtch := range reqMtchs {
			if !reqMtch.Match(req) {
				find = false
				break
			}
		}

		if find {
			return resBld
		}

	}

	return nil
}

package grpc

import (
	"fmt"
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

func NewService(
	logger ports.Logger,
	m ports.ResMtchMap[GRPCRequest],
) ports.Service[GRPCRequest] {
	return &service{
		logger,
		m,
	}
}

type service struct {
	logger ports.Logger
	m      ports.ResMtchMap[GRPCRequest]
}

func (inst *service) Handle(req GRPCRequest) (*model.SetResponse, error) {
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

func (inst *service) findMatches(req GRPCRequest) ports.ResponseBuilder[GRPCRequest] {
	for responser, mathcer := range inst.m {
		if mathcer.Match(req) {
			return responser
		}

	}

	return nil
}

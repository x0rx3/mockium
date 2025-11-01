package http

import (
	"fmt"
	"mockium/pkg/model"
	"mockium/pkg/ports"
)

func NewService(
	m ports.ResMtchMap[HTTPRequest],
) ports.Service[HTTPRequest] {
	return &service{
		m,
	}
}

type service struct {
	m ports.ResMtchMap[HTTPRequest]
}

func (inst *service) Handle(req HTTPRequest) (*model.SetResponse, error) {
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

func (inst *service) findMatches(req HTTPRequest) ports.ResponseBuilder[HTTPRequest] {
	for responser, mathcer := range inst.m {
		if mathcer.Match(req) {
			return responser
		}

	}

	return nil
}

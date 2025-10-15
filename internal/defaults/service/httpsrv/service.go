package httpsrv

import (
	"fmt"
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func NewService(
	m core.ResMtchMap[httpgw.HTTPRequest],
) core.Service[httpgw.HTTPRequest] {
	return &service{
		m,
	}
}

type service struct {
	m core.ResMtchMap[httpgw.HTTPRequest]
}

func (inst *service) Handle(req httpgw.HTTPRequest) (*model.SetResponse, error) {
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

func (inst *service) findMatches(req httpgw.HTTPRequest) core.ResponseBuilder[httpgw.HTTPRequest] {
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

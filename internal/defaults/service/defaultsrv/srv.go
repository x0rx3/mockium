package defaultsrv

import (
	"fmt"
	"mockium/pkg/core"
	"mockium/pkg/core/model"

	"go.uber.org/zap"
)

func NewService(
	log *zap.Logger,
	m core.ResMtchMap[core.Request],
) core.Service[core.Request] {
	return &service[core.Request]{
		log,
		m,
	}
}

type service[T core.Request] struct {
	log *zap.Logger
	m   core.ResMtchMap[T]
}

func (inst *service[T]) Handle(req T) (*model.SetResponse, error) {
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

func (inst *service[T]) findMatches(req T) core.ResponseBuilder[T] {
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

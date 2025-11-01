package grpc

import (
	"mockium/pkg/model"
	"mockium/pkg/ports"
	"net/http"
)

func NewHandler(
	service string,
	unarySrv map[string]ports.Service[GRPCRequest],
	steamSrv map[string]ports.Service[GRPCRequest],
) GRPCHandler[GRPCRequest] {
	return &handler{
		service:  service,
		unarySrv: unarySrv,
		steamSrv: steamSrv,
	}
}

type handler struct {
	service  string
	unarySrv map[string]ports.Service[GRPCRequest]
	steamSrv map[string]ports.Service[GRPCRequest]
}

func (inst *handler) Service() string { return inst.service }

func (inst *handler) UnaryMethods() []string {
	methods := make([]string, 0)
	for m := range inst.unarySrv {
		methods = append(methods, m)
	}

	return methods
}

func (inst *handler) StreamMethods() []string {
	methods := make([]string, 0)
	for m := range inst.steamSrv {
		methods = append(methods, m)
	}
	return methods
}

func (inst *handler) Handle(req GRPCRequest) (*model.SetResponse, error) {
	switch req.Type() {
	case "unary":
		if srv, ok := inst.unarySrv[req.Method()]; ok {
			return srv.Handle(req)
		}
	case "stream":
		if srv, ok := inst.steamSrv[req.Method()]; ok {
			return srv.Handle(req)
		}
	}

	return &model.SetResponse{
		SetStatus: http.StatusNotFound,
	}, nil
}

package httpgw

import (
	"mockium/pkg/core"
	"mockium/pkg/core/model"
	"net/http"
)

func NewHandler(
	path string,
	mthSrv map[string]core.Service[HTTPRequest],
) HTTPHandler[HTTPRequest] {
	return &handler{
		path:   path,
		mthSrv: mthSrv,
	}
}

type handler struct {
	path   string
	mthSrv map[string]core.Service[HTTPRequest]
}

func (inst *handler) Path() string { return inst.path }

func (inst *handler) SupportedMethod() []string {
	methods := make([]string, 0)
	for method := range inst.mthSrv {
		methods = append(methods, method)
	}
	return methods
}

func (inst *handler) Handle(req HTTPRequest) (*model.SetResponse, error) {
	if srv, ok := inst.mthSrv[req.Method()]; ok {
		return srv.Handle(req)
	}

	return &model.SetResponse{
		SetStatus: http.StatusNotFound,
	}, nil
}

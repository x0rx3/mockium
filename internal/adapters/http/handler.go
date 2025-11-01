package http

import (
	"mockium/internal/adapters/common"
	"mockium/pkg/model"
	"mockium/pkg/ports"
	"net/http"
)

var BuildHandler = func(logger ports.Logger, path string, handles []model.Handle) HTTPHandler[HTTPRequest] {
	mthSrv := make(map[string]ports.Service[HTTPRequest], 0)
	comparer := common.NewValueComparer()

	for _, handle := range handles {
		mthSrv[handle.MatchRequest.MustMethod] = NewService(
			ports.ResMtchMap[HTTPRequest]{
				NewResponseBuilder(handle.SetResponse): NewRequestMathcer(
					logger,
					handle,
					comparer,
				),
			},
		)
	}

	return NewHandler(path, mthSrv)
}

func NewHandler(
	path string,
	mthSrv map[string]ports.Service[HTTPRequest],
) HTTPHandler[HTTPRequest] {
	return &handler{
		path:   path,
		mthSrv: mthSrv,
	}
}

type handler struct {
	path   string
	mthSrv map[string]ports.Service[HTTPRequest]
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

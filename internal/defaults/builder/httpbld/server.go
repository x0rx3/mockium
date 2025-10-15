package httpbld

import (
	"mockium/internal/gateway/httpgw"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
)

func NewServerBuilder(
	config *core.ServerConfig,
) core.ServerBuilder {
	return &serverBuilder{
		config: config,
	}
}

type serverBuilder struct {
	config *core.ServerConfig
}

func (inst *serverBuilder) Build(logger core.Logger, templs []model.Template[[]model.Handle]) core.Server {
	handlers := make([]httpgw.HTTPHandler[httpgw.HTTPRequest], 0)

	for _, templ := range templs {
		handlers = append(
			handlers,
			BuildHandler(logger, templ.Path, templ.Handle),
		)
	}

	return httpgw.NewServer(logger, inst.config.IP+":"+string(rune(inst.config.Port)), handlers...)
}

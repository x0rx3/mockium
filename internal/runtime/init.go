package runtime

import (
	"fmt"
	"mockium/internal/defaults/builder/grpcbld"
	"mockium/internal/defaults/builder/httpbld"
	"mockium/internal/defaults/logger"
	"mockium/internal/runtime/di"
	"mockium/pkg/core"
	"mockium/pkg/core/model"
	"os"
)

var App *di.Container

func Init(config *core.Config) {
	App = di.NewContainer()

	defaultLogger, err := logger.NewZapLogger(config.DefaultLogger.Level, config.DefaultLogger.FilePath)
	if err != nil {
		fmt.Printf("failed init program: %s", err.Error())
		os.Exit(1)
	}
	App.AddLogger("default", defaultLogger)

	if httpCfg, ok := config.Servers["http"]; ok {
		builder := httpbld.NewServerBuilder(&httpCfg)
		App.AddBuilder("http", builder)

		if httpCfg.Logger.Type == "default" {
			App.AddLogger("http", defaultLogger)
		}
	}

	if grpcCfg, ok := config.Servers["grpc"]; ok {
		builder := grpcbld.NewServerBuilder(&grpcCfg)
		App.AddBuilder("grpc", builder)

		if grpcCfg.Logger.Type == "default" {
			App.AddLogger("grpc", defaultLogger)
		}
	}

}

var ParseTemplatesHandle = func(templates []model.Template[[]model.HandleTemplate]) (*model.SortedTemplates, error) {
	sortedTemplates := model.NewSortedTemplates()
	for _, template := range templates {
		newTemplate := model.Template[[]model.Handle]{}
		for _, handleTemplate := range template.Handle {
			cookies := make([]model.Cookie, 0)
			for _, cookie := range handleTemplate.MatchRequestTemplate.MustHeaders.Cookie {
				cookies = append(
					cookies,
					model.Cookie{
						Name:  cookie.Name,
						Value: cookie.Value,
					},
				)
			}

			newTemplate.Handle = append(
				newTemplate.Handle,
				model.Handle{
					MatchRequest: model.MatchRequest{
						MustType:            handleTemplate.MatchRequestTemplate.MustType,
						MustMethod:          handleTemplate.MatchRequestTemplate.MustMethod,
						MustPathParameters:  handleTemplate.MatchRequestTemplate.MustPathParameters,
						MustQueryParameters: handleTemplate.MatchRequestTemplate.MustQueryParameters,
						MustBody:            handleTemplate.MatchRequestTemplate.MustBody,
						MustHeaders: model.MustHeader{
							Cookie: cookies,
							Other:  handleTemplate.MatchRequestTemplate.MustHeaders.Other,
						},
					},
				},
			)
		}
		sortedTemplates.Add(newTemplate)
	}

	return sortedTemplates, nil
}

var ParseTemplateHandle = func(template model.Template[[]model.HandleTemplate]) (model.Template[[]model.Handle], error) {
	newTemplates := model.Template[[]model.Handle]{}

	newTemplate := model.Template[[]model.Handle]{}
	for _, handleTemplate := range template.Handle {
		cookies := make([]model.Cookie, 0)
		for _, cookie := range handleTemplate.MatchRequestTemplate.MustHeaders.Cookie {
			cookies = append(
				cookies,
				model.Cookie{
					Name:  cookie.Name,
					Value: cookie.Value,
				},
			)
		}

		newTemplate.Handle = append(
			newTemplate.Handle,
			model.Handle{
				MatchRequest: model.MatchRequest{
					MustType:            handleTemplate.MatchRequestTemplate.MustType,
					MustMethod:          handleTemplate.MatchRequestTemplate.MustMethod,
					MustPathParameters:  handleTemplate.MatchRequestTemplate.MustPathParameters,
					MustQueryParameters: handleTemplate.MatchRequestTemplate.MustQueryParameters,
					MustBody:            handleTemplate.MatchRequestTemplate.MustBody,
					MustHeaders: model.MustHeader{
						Cookie: cookies,
						Other:  handleTemplate.MatchRequestTemplate.MustHeaders.Other,
					},
				},
			},
		)
	}

	return newTemplates, nil
}

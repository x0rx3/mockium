package core

import (
	"fmt"
	"mockium/internal/adapters/grpc"
	"mockium/internal/adapters/http"
	"mockium/internal/adapters/logger"
	"mockium/pkg/model"
	"mockium/pkg/ports"
	"os"
)

// Core is the global application core instance
var Core *ApplicationCore

// Init TODO: multiply default logger
// initializes the application core with builders, validators, and loggers based on the provided config.
func Init(config *ports.Config) {
	Core = NewApplicationCore()

	defaultLogger, err := logger.NewZapLogger(config.DefaultLogger.Level, config.DefaultLogger.FilePath)
	if err != nil {
		fmt.Printf("failed init program: %s", err.Error())
		os.Exit(1)
	}

	Core.AddLogger("default", defaultLogger)

	if httpCfg, ok := config.Servers["http"]; ok {
		builder := http.NewServerBuilder(&httpCfg)
		Core.AddBuilder("http", builder)

		validator := http.NewValidator()
		Core.AddValidator("http", validator)

		if httpCfg.Logger == "default" {
			Core.AddLogger("http", defaultLogger)
		}
	}

	if grpcCfg, ok := config.Servers["grpc"]; ok {
		builder := grpc.NewServerBuilder(&grpcCfg)
		Core.AddBuilder("grpc", builder)

		if grpcCfg.Logger == "default" {
			Core.AddLogger("grpc", defaultLogger)
		}
	}

}

// ParseTemplatesHandle converts templates with HandleTemplate to templates with Handle.
var ParseTemplatesHandle = func(templates []model.Template[[]model.HandleTemplate]) (*model.SortedTemplates, error) {
	sortedTemplates := model.NewSortedTemplates()
	for _, template := range templates {
		newTemplate := model.Template[[]model.Handle]{}
		for _, handleTemplate := range template.Handle {

			var mustHeaders *model.MustHeader
			if handleTemplate.MatchRequestTemplate.MustHeaders != nil {
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
				mustHeaders.Other = handleTemplate.MatchRequestTemplate.MustHeaders.Other
				mustHeaders.Cookie = cookies
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
						MustHeaders:         mustHeaders,
					},
					SetResponse: model.SetResponse{
						SetStatus:  handleTemplate.SetResponseTemplate.SetStatus,
						SetHeaders: handleTemplate.SetResponseTemplate.SetHeaders,
						SetBody:    handleTemplate.SetResponseTemplate.SetBody,
						SetFile:    handleTemplate.SetResponseTemplate.SetFile,
					},
				},
			)
		}

		newTemplate.Name = template.Name
		newTemplate.Path = template.Path
		newTemplate.Service = template.Service
		newTemplate.Type = template.Type
		sortedTemplates.Add(newTemplate)
	}

	return sortedTemplates, nil
}

// ParseTemplateHandle converts a template with HandleTemplate to a template with Handle.
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
					MustHeaders: &model.MustHeader{
						Cookie: cookies,
						Other:  handleTemplate.MatchRequestTemplate.MustHeaders.Other,
					},
				},
			},
		)
	}

	return newTemplates, nil
}

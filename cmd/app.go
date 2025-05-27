package main

import (
	"flag"
	"mockium/internal/logging"
	"mockium/internal/service/builder"
	"mockium/internal/transport"
	"mockium/internal/transport/server"

	"go.uber.org/zap"
)

func main() {
	templateDir := flag.String("template", "../templates", "location directory with template file, default './templates'")
	address := flag.String("address", ":5000", "address with port, default ':5000'")
	flag.Parse()

	log := logging.New("debug")

	templates, err := builder.NewTemplateBuilder(log).Build(*templateDir)
	if err != nil {
		log.Error("build template", zap.Error(err))
		return
	}

	routes := make([]transport.Router, 0)
	for _, template := range templates {
		routes = append(routes, builder.BuildRoutes(log, &template))
	}

	if err := server.New(log, routes...).Start(*address); err != nil {
		log.Error("start server", zap.Error(err))
		return
	}
}

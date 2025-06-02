package main

import (
	"flag"
	"fmt"
	"mockium/internal/logging"
	"mockium/internal/service/builder"
	"mockium/internal/transport"
	"mockium/internal/transport/server"
	"os"

	"go.uber.org/zap"
)

func main() {
	templateDir := flag.String("template", "templates", "location directory with template file, default './templates'")
	address := flag.String("address", ":5000", "address with port, default ':5000'")
	logLevel := flag.String("log-level", "info", "usage log level, default 'info'")
	processLogPath := flag.String("log-dir", "log", "log direcrectory, default 'log'")
	flag.Parse()

	log, err := logging.NewZapLogger(*logLevel, *processLogPath)
	if err != nil {
		fmt.Printf("failed init logger: %s", err.Error())
		os.Exit(1)
	}

	procLogger, err := logging.NewProcessLogger(log, *processLogPath, "requests", 10)
	if err != nil {
		log.Error("init process logger", zap.Error(err))
		os.Exit(1)
	}
	defer procLogger.Close()

	templates, err := builder.NewTemplateBuilder(log).Build(*templateDir)
	if err != nil {
		log.Error("build template", zap.Error(err))
		os.Exit(1)
	}

	routes := make([]transport.Router, 0)
	for _, template := range templates {
		routes = append(routes, builder.BuildRoutes(log, procLogger, &template))
	}

	if err := server.New(log, routes...).Start(*address); err != nil {
		log.Error("start server", zap.Error(err))
		os.Exit(1)
	}
}

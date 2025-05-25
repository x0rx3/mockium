package main

import (
	"flag"
	"gomock/internal/logging"
	"gomock/internal/service/builder"
	"gomock/internal/transport"
	"gomock/internal/transport/server"

	"go.uber.org/zap"
)

func main() {
	templateDir := flag.String("dir", "templtes", "location directory with template file, default './templates'")
	address := flag.String("addr", "127.0.0.1:5000", "address with port, default '127.0.0.1:5000'")
	flag.Parse()

	log := logging.New("debug")

	templates, err := builder.NewTemplateBuilder(log).Build(*templateDir)
	if err != nil {
		log.Error("error build template", zap.Error(err))
		return
	}

	r := make([]transport.Router, 0)
	for _, template := range templates {
		r = append(r, builder.BuildRoutes(log, &template)...)
	}

	if err := server.New(log, r...).Start(*address); err != nil {
		log.Error("error start server", zap.Error(err))
		return
	}
}

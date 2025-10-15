package main

import (
	"flag"
	"mockium/internal/io"
	"mockium/internal/runtime"
	"os"
)

func main() {
	templateDir := flag.String("template", "templates", "location directory with template file, default './templates'")
	config := flag.String("config", "config.yaml", "location file with config, default './config.yaml'")
	flag.Parse()

	runtime.Init(io.LoadConfig(*config))

	receivedTemplates, err := io.LoadTemplates(*templateDir)
	if err != nil {
		runtime.App.GetLogger("default").Error("failed load templates", err)
		os.Exit(1)
	}

	if err := runtime.App.ValidateTemplates(receivedTemplates); err != nil {
		runtime.App.GetLogger("default").Error("failed load templates", err)
		os.Exit(1)
	}

	internalTemplate, err := runtime.ParseTemplatesHandle(receivedTemplates)
	if err != nil {
		runtime.App.GetLogger("default").Error("failed parse templates", err)
		os.Exit(1)
	}

	if err := runtime.App.BuildServers(internalTemplate); err != nil {
		runtime.App.GetLogger("default").Error("failed build servers", err)
		os.Exit(1)
	}

	runtime.App.StartAllServers()
	defer runtime.App.StopAllServers()
}

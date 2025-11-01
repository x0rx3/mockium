package main

import (
	"flag"
	"fmt"
	"mockium/internal/core"
	"mockium/internal/io"
	"os"
)

func main() {
	templateDir := flag.String("template", "templates", "location directory with template file, default './templates'")
	configPath := flag.String("config", "config.yaml", "location file with config, default './config.yaml'")
	flag.Parse()

	appConfig, err := io.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("failed load application config: %s", err.Error())
		os.Exit(1)
	}

	core.Init(appConfig)

	receivedTemplates, err := io.LoadTemplates(*templateDir)
	if err != nil {
		core.Core.GetLogger("default").Error("failed load templates", err)
		os.Exit(1)
	}

	if err := core.Core.ValidateTemplates(receivedTemplates); err != nil {
		core.Core.GetLogger("default").Error("failed load templates", err)
		os.Exit(1)
	}

	internalTemplate, err := core.ParseTemplatesHandle(receivedTemplates)
	if err != nil {
		core.Core.GetLogger("default").Error("failed parse templates", err)
		os.Exit(1)
	}

	if err := core.Core.BuildServers(internalTemplate); err != nil {
		core.Core.GetLogger("default").Error("failed build servers", err)
		os.Exit(1)
	}

	core.Core.StartAllServers()
	defer core.Core.StopAllServers()

	select {}
}

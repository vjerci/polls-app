package main

import (
	"fmt"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/app"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/config"
)

func main() {
	if err := config.Setup(); err != nil {
		panic(fmt.Errorf("failed to setup config %w", err))
	}

	settings := config.Get()

	server, err := app.New(settings)
	if err != nil {
		panic(fmt.Errorf("failed to build server %w", err))
	}

	err = server.Start(settings.HTTPPort)
	if err != nil {
		panic(fmt.Errorf("server running error  %w", err))
	}
}

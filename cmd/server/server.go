package main

import (
	"fmt"
	"time"

	"github.com/vjerci/polls-app/pkg/app"
	"github.com/vjerci/polls-app/pkg/config"
)

// wait for postgre to come online.
var startupDelay = 5 * time.Second

func main() {
	if err := config.Setup(); err != nil {
		panic(fmt.Errorf("failed to setup config %w", err))
	}

	settings := config.Get()

	time.Sleep(startupDelay)

	server, err := app.New(settings)
	if err != nil {
		panic(fmt.Errorf("failed to build server %w", err))
	}

	err = server.Start(":" + settings.HTTPPort)
	if err != nil {
		panic(fmt.Errorf("server running error  %w", err))
	}
}

package main

import (
	"fmt"
	"log"
	"net"
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

	lis, err := net.Listen("tcp", ":"+settings.GRPCPort)
	if err != nil {
		panic(fmt.Errorf("failed to listen for grpc on port %s", settings.GRPCPort))
	}

	server, err := app.NewGrpc(settings)
	if err != nil {
		panic(fmt.Errorf("failed to build server %w", err))
	}

	log.Printf("grpc server listening on port %s", settings.GRPCPort)

	err = server.Serve(lis)
	if err != nil {
		panic(fmt.Errorf("server running error  %w", err))
	}
}

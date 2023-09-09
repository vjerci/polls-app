package config

import (
	"errors"
	"os"
)

var ErrPostgresURLEmpty = errors.New("postgres url is not specified")
var ErrJWTKeyEmpty = errors.New("jwt private key is not specified")
var ErrHTTPPortEmpty = errors.New("http port is not specified")
var ErrGRPCPortEmpty = errors.New("grpc port is not specified")

var config Config

type Config struct {
	HTTPPort      string
	GRPCPort      string
	PostgresURL   string
	JWTSigningKey string
	Debug         bool
}

func Setup() error {
	postgres := os.Getenv("POSTGRES_URL")
	if postgres == "" {
		return ErrPostgresURLEmpty
	}

	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	jwt := os.Getenv("JWT_KEY")
	if jwt == "" {
		return ErrJWTKeyEmpty
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return ErrHTTPPortEmpty
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		return ErrGRPCPortEmpty
	}

	config = Config{
		PostgresURL:   postgres,
		JWTSigningKey: jwt,
		Debug:         debug,
		HTTPPort:      ":" + httpPort,
		GRPCPort:      grpcPort,
	}

	return nil
}

func Get() Config {
	return config
}

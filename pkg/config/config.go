package config

import (
	"errors"
	"os"
)

var ErrPostgresURLEmpty = errors.New("postgres url is not specified")
var ErrJWTKeyEmpty = errors.New("jwt private key is not specified")
var ErrPortEmpty = errors.New("port is not specified")

var config Config

type Config struct {
	Port          string
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

	port := os.Getenv("PORT")
	if port == "" {
		return ErrPortEmpty
	}

	config = Config{
		PostgresURL:   postgres,
		JWTSigningKey: jwt,
		Debug:         debug,
		Port:          ":" + port,
	}

	return nil
}

func Get() Config {
	return config
}

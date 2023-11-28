package config

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrReadConfig = errors.New("failed to read config")
var ErrUnmarshalConfig = errors.New("failed to unmarshal config")

var ErrPostgresURLEmpty = errors.New("postgres url is not specified")
var ErrJWTKeyEmpty = errors.New("jwt private key is not specified")
var ErrHTTPPortEmpty = errors.New("http port is not specified")
var ErrGoogleClientIDEmpty = errors.New("google client id empty")

var config Config

type Config struct {
	PostgresURL string `mapstructure:"POSTGRES_URL"`

	HTTPPort string `mapstructure:"HTTP_PORT"`
	GRPCPort string `mapstructure:"GRPC_PORT"`

	JWTSigningKey  string `mapstructure:"JWT_KEY"`
	GoogleClientID string `mapstructure:"GOOGLE_CLIENT_ID"`

	Debug bool `mapstructure:"DEBUG"`
}

func Setup() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return errors.Join(ErrReadConfig, err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return errors.Join(ErrUnmarshalConfig, err)
	}

	if config.PostgresURL == "" {
		return ErrPostgresURLEmpty
	}

	if config.JWTSigningKey == "" {
		return ErrJWTKeyEmpty
	}

	if config.HTTPPort == "" {
		return ErrHTTPPortEmpty
	}

	if config.GoogleClientID == "" {
		return ErrGoogleClientIDEmpty
	}

	return nil
}

func Get() Config {
	return config
}

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConnect = errors.New("unable to connect to postgres")
var ErrPing = errors.New("unable to ping postgres")

type DB struct {
	PostgresURL string
	Pool        *pgxpool.Pool
}

// postgresURL := "postgres://username:password@localhost:5432/database_name"
func New(postgresURL string) *DB {
	return &DB{
		PostgresURL: postgresURL,
		Pool:        nil,
	}
}

func (client *DB) Connect() error {
	conn, err := pgxpool.New(context.Background(), client.PostgresURL)
	if err != nil {
		return errors.Join(ErrConnect, err)
	}

	client.Pool = conn

	if err = client.Pool.Ping(context.Background()); err != nil {
		return errors.Join(ErrPing, err)
	}

	return nil
}

func (client *DB) Close() {
	client.Pool.Close()
}

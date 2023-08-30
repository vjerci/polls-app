package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConnect = errors.New("unable to connect to postgres")
var ErrPing = errors.New("unable to ping postgres")

type Client struct {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	PostgresURL string
	Pool        *pgxpool.Pool
}

func New(postgresURL string) *Client {
	return &Client{
		PostgresURL: postgresURL,
	}
}

func (client *Client) Connect() error {
	conn, err := pgxpool.New(context.Background(), client.PostgresURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConnect, err)
	}

	client.Pool = conn

	if err = client.Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("%w: %w", ErrPing, err)
	}

	return nil
}

func (client *Client) Close() {
	client.Pool.Close()
}

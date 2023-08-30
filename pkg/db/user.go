package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrGetUserNoRows = errors.New("couldn't find user for a given user_id, no rows returned")
var ErrGetUserQueryFailed = errors.New("execution of get user query failed")

func (client *Client) GetUser(userID string) (name string, err error) {
	err = client.Pool.QueryRow(context.Background(), "SELECT name FROM users WHERE id=$1", userID).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrGetUserNoRows
		}

		return "", errors.Join(ErrGetUserQueryFailed, err)
	}

	return name, nil
}

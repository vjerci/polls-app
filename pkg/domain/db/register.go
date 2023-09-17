package db

import (
	"context"
	"errors"
	"strings"
)

var ErrCreateUserExecuteQuery = errors.New("error while executing register query")
var ErrCreateUserInsertCount = errors.New("error while inserting user, 0 rows were affected due to duplication")

func (client *DB) CreateUser(userID string, name string) error {
	_, err := client.Pool.Exec(context.Background(), "INSERT INTO users(id, name) VALUES ($1, $2);", userID, name)
	if err != nil {
		if strings.Index(err.Error(), "ERROR: duplicate key value violates unique constraint") == 0 {
			return ErrCreateUserInsertCount
		}

		return errors.Join(ErrCreateUserExecuteQuery, err)
	}

	return nil
}

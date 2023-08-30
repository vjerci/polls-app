package db

import (
	"context"
	"errors"
	"fmt"
)

var ErrRegisterExecuteQuery = errors.New("error while executing register query")
var ErrRegisterInsertCount = errors.New("error while inserting user, 0 rows were affected")

func (client *Client) CreateUser(userID string, groupID string, name string) error {
	res, err := client.Pool.Exec(context.Background(), "INSERT INTO users(id, group_id, name) VALUES ($1, $2, $3);", userID, groupID, name)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRegisterExecuteQuery, err)
	}

	if res.RowsAffected() == 0 {
		return ErrRegisterInsertCount
	}

	return nil
}

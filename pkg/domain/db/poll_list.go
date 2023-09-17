package db

import (
	"context"
	"errors"

	pgx "github.com/jackc/pgx/v5"
)

var PollListLimit = 10

type PollListData struct {
	Name string `db:"name"`
	ID   string `db:"id"`
}

var ErrPollListQuery = errors.New("errors querying poll list")
var ErrPollListCollectRows = errors.New("error collecting poll list rows")

func (client *DB) GetPollList(page int) ([]PollListData, error) {
	rows, err := client.Pool.Query(
		context.Background(),
		"SELECT name, id FROM polls ORDER BY date_created OFFSET $1 LIMIT $2;",
		page*PollListLimit,
		PollListLimit,
	)
	if err != nil {
		return nil, errors.Join(ErrPollListQuery, err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[PollListData])
	if err != nil {
		return nil, errors.Join(ErrPollListCollectRows, err)
	}

	return data, nil
}

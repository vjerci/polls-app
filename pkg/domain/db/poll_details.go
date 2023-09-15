package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type PollDetailsResponse struct {
	Name string
	ID   string
}

var ErrPollDetailsQuery = errors.New("failed to query polls table")
var ErrPollDetailsNotFound = errors.New("poll does not exists")

func (client *Client) GetPollDetails(pollID string) (*PollDetailsResponse, error) {
	var pollName string
	err := client.Pool.QueryRow(
		context.Background(),
		"SELECT name FROM polls WHERE id=$1;",
		pollID,
	).Scan(&pollName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(ErrPollDetailsNotFound, err)
		}

		return nil, errors.Join(ErrPollListQuery, err)
	}

	return &PollDetailsResponse{
		Name: pollName,
		ID:   pollID,
	}, nil
}

type PollDetailsAnswer struct {
	Name  string `db:"name"`
	ID    string `db:"id"`
	Count int    `db:"count"`
}

var ErrPollDetailsAnswerQuery = errors.New("failed to query answers table")
var ErrPollDetailsAnswerScan = errors.New("failed to scan answers table")

func (client *Client) GetPollDetailsAnswers(pollID string) ([]PollDetailsAnswer, error) {
	rows, err := client.Pool.Query(
		context.Background(),
		`SELECT answers.name AS name, answers.id AS id, count(user_answers.id) as count 
			FROM answers
		 JOIN user_answers ON  user_answers.answer_id=answers.id
			GROUP_BY user_answers.answer_id;
		WHERE answers.poll_id=$1 
		`,
		pollID,
	)
	if err != nil {
		return nil, errors.Join(ErrPollDetailsAnswerQuery, err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[PollDetailsAnswer])
	if err != nil {
		return nil, errors.Join(ErrPollDetailsAnswerScan, err)
	}

	return data, nil
}

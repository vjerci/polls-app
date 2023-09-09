package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type CreatePollResponse struct {
	Name    string
	ID      string
	Answers []CreatePollAnswer
}

type CreatePollAnswer struct {
	Name string
	ID   string
}

var ErrCreatePollInsert = errors.New("failed to insert entry into polls table")
var ErrCreatePollAnswer = errors.New("failed to insert entry into answer table")

func (client *Client) CreatePoll(name string, answers []string) (*CreatePollResponse, error) {
	var pollID string

	err := client.Pool.QueryRow(
		context.Background(),
		`INSERT INTO polls (name, date_created) VALUES
		($1, $2) RETURNING id;`,
		name,
		time.Now().Unix(),
	).Scan(&pollID)
	if err != nil {
		return nil, errors.Join(ErrCreatePollInsert, err)
	}

	batch := &pgx.Batch{}

	for _, answer := range answers {
		batch.Queue(`INSERT INTO poll_answers (poll_id, name) VALUES ($1, $2) RETURNING id`, pollID, answer)
	}

	createResult := &CreatePollResponse{
		Name:    name,
		ID:      pollID,
		Answers: nil,
	}

	batchResult := client.Pool.SendBatch(context.Background(), batch)
	defer batchResult.Close()

	for answerPointer := 0; answerPointer < batch.Len(); answerPointer++ {
		var answerID string

		err := batchResult.QueryRow().Scan(&answerID)
		if err != nil {
			return nil, errors.Join(ErrCreatePollAnswer, err)
		}

		createResult.Answers = append(createResult.Answers, CreatePollAnswer{
			Name: answers[answerPointer],
			ID:   answerID,
		})
	}

	return createResult, nil
}

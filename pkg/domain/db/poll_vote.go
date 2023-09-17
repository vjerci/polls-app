package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrExistsPollAnswerNotFound = errors.New("poll answer with given answer id and poll id does not exist")
var ErrExistsPollAnswerQuery = errors.New("failed to query answers")

func (client *Client) ExistsPollAnswer(pollID, answerID string) error {
	_, err := client.Pool.Query(
		context.Background(),
		"SELECT name FROM answers WHERE id=$1 AND poll_id=$2;",
		answerID,
		pollID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(ErrExistsPollAnswerNotFound, err)
		}

		return errors.Join(ErrExistsPollAnswerQuery, err)
	}

	return nil
}

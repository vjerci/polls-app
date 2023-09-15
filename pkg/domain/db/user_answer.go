package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrUserAnswerNotFound = errors.New("no record for user answer")
var ErrUserAnswerQuery = errors.New("failed to query user answer")

func (client *Client) GetUserAnswer(pollID, userID string) (answerID string, err error) {
	err = client.Pool.QueryRow(
		context.Background(),
		`SELECT answers.id AS id FROM answers
		 JOIN user_answers ON  user_answers.answer_id=answers.id
		 WHERE answers.poll_id=$1 AND user_answers.user_id= $2
		`,
		pollID,
		userID,
	).Scan(&answerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.Join(ErrUserAnswerNotFound, err)
		}

		return "", errors.Join(ErrPollListQuery, err)
	}

	return answerID, nil
}

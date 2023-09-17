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

var ErrUpsertUserAnswerQuery = errors.New("failed to query user answers for upserting")
var ErrUpsertUserAnswerUpdate = errors.New("failed to update user answer for upserting")
var ErrUpsertUserAnswerInsert = errors.New("failed to insert user answer for upserting")

func (client *Client) UpsertUserAnswer(pollID, answerID, userID string) (modified bool, err error) {
	oldAnswerID, err := client.GetUserAnswer(pollID, userID)

	if err != nil && !errors.Is(err, ErrUserAnswerNotFound) {
		return false, errors.Join(ErrUpsertUserAnswerQuery, err)
	}

	if oldAnswerID != "" {
		_, err := client.Pool.Query(
			context.Background(),
			"UPDATE user_answers SET answer_id=$1 WHERE answer_id=$1 AND user_id=$2;",
			answerID,
			oldAnswerID,
			userID,
		)

		if err != nil {
			return false, errors.Join(ErrUpsertUserAnswerUpdate, err)
		}

		return true, nil
	}

	var insertedUserID string
	err = client.Pool.QueryRow(
		context.Background(),
		`INSERT INTO user_answers (answer_id, user_id) VALUES
		($1, $2) RETURNING user_id`,
		answerID,
		userID,
	).Scan(&insertedUserID)

	if err != nil {
		return false, errors.Join(ErrUpsertUserAnswerInsert, err)
	}

	return false, nil
}

package model

import (
	"errors"
)

type PollVoteRequest struct {
	PollID   string
	AnswerID string
	UserID   string
}

type PollVoteResponse struct {
	ModifiedAnswer bool
}

type PollVoteRepository interface {
	ExistsPollAnswer(pollID string, answerID string) (err error)
	UpsertUserAnswer(pollID, answerID, userID string) (modified bool, err error)
}

var ErrPollVoteUserIDEmpty = errors.New("failed to submit vote, userID is empty")
var ErrPollVotePollIDEmpty = errors.New("failed to submit vote, pollID is empty")
var ErrPollVoteAnswerIDEmpty = errors.New("failed to submit vote, answerID is empty")

var ErrPollVoteAnswerNotFound = errors.New("failed to submit vote, combination of pollID and answerID doesn't exist")
var ErrPollVoteUpsertFailed = errors.New("failed to submit vote, upsert of user answer failed")

func (client *Client) PollVote(data *PollVoteRequest) (*PollVoteResponse, error) {
	if data.PollID == "" {
		return nil, ErrPollVotePollIDEmpty
	}

	if data.AnswerID == "" {
		return nil, ErrPollVoteAnswerIDEmpty
	}

	if data.UserID == "" {
		return nil, ErrPollVoteUserIDEmpty
	}

	err := client.PollVoteDB.ExistsPollAnswer(data.PollID, data.AnswerID)
	if err != nil {
		return nil, errors.Join(ErrPollVoteAnswerNotFound, err)
	}

	modified, err := client.PollVoteDB.UpsertUserAnswer(data.PollID, data.AnswerID, data.UserID)
	if err != nil {
		return nil, errors.Join(ErrPollVoteUpsertFailed, err)
	}

	return &PollVoteResponse{
		ModifiedAnswer: modified,
	}, nil
}

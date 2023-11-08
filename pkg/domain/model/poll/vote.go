package poll

import (
	"errors"
)

type VoteRequest struct {
	PollID   string
	AnswerID string
	UserID   string
}

type VoteResponse struct {
	ModifiedAnswer bool
}

type VoteRepository interface {
	ExistsPollAnswer(pollID string, answerID string) (err error)
	UpsertUserAnswer(pollID, answerID, userID string) (modified bool, err error)
}

type VoteModel struct {
	VoteDB VoteRepository
}

var ErrVoteUserIDEmpty = errors.New("failed to submit vote, userID is empty")
var ErrVoteAnswerIDEmpty = errors.New("failed to submit vote, answerID is empty")
var ErrVotePollIDEmpty = errors.New("failed to submit vote, pollID is empty")

var ErrVoteAnswerNotFound = errors.New("failed to submit vote, combination of pollID and answerID doesn't exist")
var ErrVoteUpsertFailed = errors.New("failed to submit vote, upsert of user answer failed")

func (model *VoteModel) Do(data *VoteRequest) (*VoteResponse, error) {
	if data.PollID == "" {
		return nil, ErrVotePollIDEmpty
	}

	if data.AnswerID == "" {
		return nil, ErrVoteAnswerIDEmpty
	}

	if data.UserID == "" {
		return nil, ErrVoteUserIDEmpty
	}

	err := model.VoteDB.ExistsPollAnswer(data.PollID, data.AnswerID)
	if err != nil {
		return nil, errors.Join(ErrVoteAnswerNotFound, err)
	}

	modified, err := model.VoteDB.UpsertUserAnswer(data.PollID, data.AnswerID, data.UserID)
	if err != nil {
		return nil, errors.Join(ErrVoteUpsertFailed, err)
	}

	return &VoteResponse{
		ModifiedAnswer: modified,
	}, nil
}

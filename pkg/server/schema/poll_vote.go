package schema

import (
	"errors"
	"net/http"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollVoteRequest struct {
	AnswerID string `json:"answer_id"`
}

type PollVoteResponse struct {
	ModifiedAnswer bool `json:"modified_answer"`
}

func (mapper *Map) PollVoteResponse(input *model.PollVoteResponse) *PollVoteResponse {
	return &PollVoteResponse{
		ModifiedAnswer: input.ModifiedAnswer,
	}
}

var ErrPollVoteInvalidVote = &UserVisibleError{
	Err:    model.ErrPollVoteAnswerNotFound,
	Status: http.StatusNotFound,
}

var ErrPollVoteInvalidPollID = &UserVisibleError{
	Err:    model.ErrPollVotePollIDEmpty,
	Status: http.StatusBadRequest,
}
var ErrPollVoteInvalidAnswerID = &UserVisibleError{
	Err:    model.ErrPollVoteAnswerIDEmpty,
	Status: http.StatusBadRequest,
}
var ErrPollVoteInvalidUserID = &UserVisibleError{
	Err:    model.ErrPollVoteUserIDEmpty,
	Status: http.StatusBadRequest,
}
var handledPollVoteErrors = []*UserVisibleError{
	ErrPollVoteInvalidVote,
	ErrPollVoteInvalidPollID,
	ErrPollVoteInvalidAnswerID,
	ErrPollVoteInvalidUserID,
}

var ErrPollVoteJSONDecode = &UserVisibleError{
	Err:    errors.New("failed to decode poll vote json body"),
	Status: http.StatusBadRequest,
}

var ErrPollVoteModel = errors.New("model poll vote failed to cast vote")

func (mapper *Map) PollVoteError(err error) error {
	for _, targetError := range handledPollVoteErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrPollVoteModel, err)
}

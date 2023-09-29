package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollVoteRequest struct {
	AnswerID string `json:"answer_id"`
}

type PollVoteSchemaMap struct{}

type PollVoteResponse struct {
	ModifiedAnswer bool `json:"modified_answer"`
}

func (mapper *PollVoteSchemaMap) Response(input *model.PollVoteResponse) *PollVoteResponse {
	return &PollVoteResponse{
		ModifiedAnswer: input.ModifiedAnswer,
	}
}

var ErrPollVoteInvalidVote = &echo.HTTPError{
	Message:  model.ErrPollVoteAnswerNotFound,
	Code:     http.StatusNotFound,
	Internal: nil,
}

var ErrPollVoteInvalidPollID = &echo.HTTPError{
	Message:  model.ErrPollVotePollIDEmpty,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollVoteInvalidAnswerID = &echo.HTTPError{
	Message:  model.ErrPollVoteAnswerIDEmpty,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollVoteInvalidUserID = &echo.HTTPError{
	Message:  model.ErrPollVoteUserIDEmpty,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollVoteJSONDecode = &echo.HTTPError{
	Message:  errors.New("failed to decode poll vote json body"),
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollVoteModel = &echo.HTTPError{
	Message:  errors.New("internal server error"),
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollVoteSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrPollVoteAnswerNotFound) {
		return ErrPollVoteInvalidVote.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollVotePollIDEmpty) {
		return ErrPollVoteInvalidPollID.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollVoteAnswerIDEmpty) {
		return ErrPollVoteInvalidAnswerID.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollVoteUserIDEmpty) {
		return ErrPollVoteInvalidUserID.WithInternal(err)
	}

	return ErrPollVoteModel.SetInternal(err)
}

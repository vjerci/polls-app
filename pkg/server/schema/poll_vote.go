package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
)

type PollVoteRequest struct {
	AnswerID string `json:"answer_id"`
}

type PollVoteSchemaMap struct{}

type PollVoteResponse struct {
	ModifiedAnswer bool `json:"modified_answer"`
}

func (mapper *PollVoteSchemaMap) Response(input *poll.VoteResponse) *PollVoteResponse {
	return &PollVoteResponse{
		ModifiedAnswer: input.ModifiedAnswer,
	}
}

var ErrPollVoteInvalidVote = &echo.HTTPError{
	Message:  `couldn't find an answer for given input "answer_id"`,
	Code:     http.StatusNotFound,
	Internal: nil,
}

var ErrPollVoteInvalidPollID = &echo.HTTPError{
	Message:  "inputed poll id is empty",
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollVoteInvalidAnswerID = &echo.HTTPError{
	Message:  `input field "answer_id" can't be empty`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollVoteInvalidUserID = &echo.HTTPError{
	Message:  `user_id can't be empty`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollVoteJSONDecode = &echo.HTTPError{
	Message:  "failed to decode poll vote json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollVoteModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollVoteSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, poll.ErrVoteAnswerNotFound) {
		return ErrPollVoteInvalidVote.WithInternal(err)
	}

	if errors.Is(err, poll.ErrVotePollIDEmpty) {
		return ErrPollVoteInvalidPollID.WithInternal(err)
	}

	if errors.Is(err, poll.ErrVoteAnswerIDEmpty) {
		return ErrPollVoteInvalidAnswerID.WithInternal(err)
	}

	if errors.Is(err, poll.ErrVoteUserIDEmpty) {
		return ErrPollVoteInvalidUserID.WithInternal(err)
	}

	return ErrPollVoteModel.SetInternal(err)
}

package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollDetailsSchemaMap struct{}

type PollDetailsResponse struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	UserVote string              `json:"user_vote"`
	Answers  []PollDetailsAnswer `json:"answers"`
}

type PollDetailsAnswer struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	VotesCount int    `json:"votes_count"`
}

func (mapper *PollDetailsSchemaMap) Response(input *model.PollDetailsResponse) *PollDetailsResponse {
	answers := make([]PollDetailsAnswer, len(input.Answers))
	for i, v := range input.Answers {
		answers[i] = PollDetailsAnswer{
			Name:       v.Name,
			ID:         v.ID,
			VotesCount: v.VotesCount,
		}
	}

	return &PollDetailsResponse{
		Name:     input.Name,
		ID:       input.ID,
		UserVote: input.UserAnswer,
		Answers:  answers,
	}
}

var ErrPollDetailsEmptyPollID = &echo.HTTPError{
	Message:  "inputted poll id can't be empty",
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollDetailsEmptyUserID = &echo.HTTPError{
	Message:  "inputted user id can't be empty",
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollDetailsNotFound = &echo.HTTPError{
	Message:  "poll with a given id not found",
	Code:     http.StatusNotFound,
	Internal: nil,
}

var ErrPollDetailsModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollDetailsSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrPollDetailsIDEmpty) {
		return ErrPollDetailsEmptyPollID.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollDetailsUserIDEmpty) {
		return ErrPollDetailsEmptyUserID.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollDetailsNoPoll) {
		return ErrPollDetailsNotFound.WithInternal(err)
	}

	return ErrPollDetailsModel.WithInternal(err)
}

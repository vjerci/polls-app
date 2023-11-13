package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
)

type PollCreateRequest struct {
	Name    string   `json:"name"`
	Answers []string `json:"answers"`
}

func (req *PollCreateRequest) ToModel() *poll.CreateRequest {
	return &poll.CreateRequest{
		Name:    req.Name,
		Answers: req.Answers,
	}
}

type PollCreateSchemaMap struct{}

type PollCreateResponse struct {
	ID         string   `json:"id"`
	AnswersIDS []string `json:"answers_ids"`
}

func (mapper *PollCreateSchemaMap) Response(input *poll.CreateResponse) *PollCreateResponse {
	return &PollCreateResponse{
		ID:         input.PollID,
		AnswersIDS: input.AnswersIDS,
	}
}

var ErrPollCreateNameEmpty = &echo.HTTPError{
	Message:  `input field "name" can't be empty`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollCreateAnswersLen = &echo.HTTPError{
	Message:  `input field "answers" needs minimum 2 answers`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollCreateAnswerEmpty = &echo.HTTPError{
	Message:  `input field "answers" can't contain empty strings`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollCreateJSONDecode = &echo.HTTPError{
	Message:  "failed to decode poll_create input json",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollCreateModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollCreateSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, poll.ErrCreateNameEmpty) {
		return ErrPollCreateNameEmpty.WithInternal(err)
	}

	if errors.Is(err, poll.ErrCreateAnswersLen) {
		return ErrPollCreateAnswersLen.WithInternal(err)
	}

	if errors.Is(err, poll.ErrCreateAnswerEmpty) {
		return ErrPollCreateAnswerEmpty.WithInternal(err)
	}

	return ErrPollCreateModel.WithInternal(ErrPollCreateModel)
}

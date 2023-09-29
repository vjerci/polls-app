package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollCreateRequest struct {
	Name    string   `json:"name"`
	Answers []string `json:"answers"`
}

func (req *PollCreateRequest) ToModel() *model.PollCreateRequest {
	return &model.PollCreateRequest{
		Name:    req.Name,
		Answers: req.Answers,
	}
}

type PollCreateSchemaMap struct{}

type PollCreateResponse struct {
	ID         string   `json:"id"`
	AnswersIDS []string `json:"answers_ids"`
}

func (mapper *PollCreateSchemaMap) Response(input *model.PollCreateResponse) *PollCreateResponse {
	return &PollCreateResponse{
		ID:         input.PollID,
		AnswersIDS: input.AnswersIDS,
	}
}

var ErrPollCreateNameEmpty = &echo.HTTPError{
	Message:  model.ErrPollCreateNameEmpty,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollCreateAnswersLen = &echo.HTTPError{
	Message:  model.ErrPollCreateAnswersLen,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollCreateAnswerEmpty = &echo.HTTPError{
	Message:  model.ErrPollCreateAnswerEmpty,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollCreateJSONDecode = &echo.HTTPError{
	Message:  errors.New("failed to decode create poll json body"),
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollCreateModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollCreateSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrPollCreateNameEmpty) {
		return ErrPollCreateNameEmpty.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollCreateAnswersLen) {
		return ErrPollCreateAnswersLen.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollCreateAnswerEmpty) {
		return ErrPollCreateAnswerEmpty.WithInternal(err)
	}

	return ErrPollCreateModel.WithInternal(ErrPollCreateModel)
}

package schema

import (
	"errors"
	"net/http"

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

var ErrPollCreateNameEmpty = &UserVisibleError{
	Err:    model.ErrPollCreateNameEmpty,
	Status: http.StatusBadRequest,
}
var ErrPollCreateAnswersLen = &UserVisibleError{
	Err:    model.ErrPollCreateAnswersLen,
	Status: http.StatusBadRequest,
}
var ErrPollCreateAnswerEmpty = &UserVisibleError{
	Err:    model.ErrPollCreateAnswerEmpty,
	Status: http.StatusBadRequest,
}
var handledPollCreateErrors = []*UserVisibleError{
	ErrPollCreateNameEmpty,
	ErrPollCreateAnswersLen,
	ErrPollCreateAnswerEmpty,
}

var ErrPollCreateJSONDecode = &UserVisibleError{
	Err:    errors.New("failed to decode create poll json body"),
	Status: http.StatusBadRequest,
}

var ErrPollCreateModel = errors.New("model failed to create poll")

func (mapper *PollCreateSchemaMap) ErrorHandler(err error) error {
	for _, targetError := range handledPollCreateErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrPollCreateModel, err)
}

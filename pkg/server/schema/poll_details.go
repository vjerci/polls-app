package schema

import (
	"errors"
	"net/http"

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

var ErrPollDetailsEmptyPollID = &UserVisibleError{
	Err:    model.ErrPollDetailsIDEmpty,
	Status: http.StatusBadRequest,
}
var ErrPollDetailsEmptyUserID = &UserVisibleError{
	Err:    model.ErrPollDetailsUserIDEmpty,
	Status: http.StatusBadRequest,
}
var ErrPollDetailsNotFound = &UserVisibleError{
	Err:    model.ErrPollDetailsNoPoll,
	Status: http.StatusNotFound,
}
var handledPollDetailsErrors = []*UserVisibleError{
	ErrPollDetailsEmptyPollID,
	ErrPollDetailsEmptyUserID,
	ErrPollDetailsNotFound,
}

var ErrPollDetailsModel = errors.New("model failed to get poll details")

func (mapper *PollDetailsSchemaMap) ErrorHandler(err error) error {
	for _, targetError := range handledPollDetailsErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrPollDetailsModel, err)
}

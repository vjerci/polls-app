package schema

import (
	"errors"
	"net/http"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollDetailsResponse struct {
	Name     string              `json:"name"`
	ID       string              `json:"id"`
	Answers  []PollDetailsAnswer `json:"answers"`
	UserVote string              `json:"user_vote"`
}

type PollDetailsAnswer struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	VotesCount int    `json:"votes_count"`
}

func (mapper *Map) PollDetailsResponse(input *model.PollDetailsResponse) *PollDetailsResponse {
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

func (mapper *Map) PollDetailsError(err error) error {
	for _, targetError := range handledPollDetailsErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrPollDetailsModel, err)
}

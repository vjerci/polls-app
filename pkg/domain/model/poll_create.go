package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
)

var CreatePollMinAnswers = 2

type PollCreateRequest struct {
	Name    string
	Answers []string
}

type PollCreateResponse struct {
	PollID     string   `json:"poll_id"`
	AnswersIDS []string `json:"answer_ids"`
}

type PollCreateRepository interface {
	CreatePoll(name string, answers []string) (*db.PollCreateResponse, error)
}

type PollCreateModel struct {
	PollCreateDB PollCreateRepository
}

var ErrPollCreateNameEmpty = errors.New("poll name cannot be empty")
var ErrPollCreateAnswersLen = errors.New("poll needs at least 2 answers")
var ErrPollCreateAnswerEmpty = errors.New("poll answer name can't be empty")
var ErrPollCreateDB = errors.New("failed to create poll in db")

func (model *PollCreateModel) Create(data *PollCreateRequest) (*PollCreateResponse, error) {
	if data.Name == "" {
		return nil, ErrPollCreateNameEmpty
	}

	if len(data.Answers) < CreatePollMinAnswers {
		return nil, ErrPollCreateAnswersLen
	}

	for _, answer := range data.Answers {
		if answer == "" {
			return nil, ErrPollCreateAnswerEmpty
		}
	}

	dbPolls, err := model.PollCreateDB.CreatePoll(data.Name, data.Answers)
	if err != nil {
		return nil, errors.Join(ErrPollCreateDB, err)
	}

	resp := &PollCreateResponse{
		PollID:     dbPolls.ID,
		AnswersIDS: nil,
	}

	for _, answer := range dbPolls.Answers {
		resp.AnswersIDS = append(resp.AnswersIDS, answer.ID)
	}

	return resp, nil
}

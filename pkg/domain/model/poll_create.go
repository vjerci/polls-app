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
	PollID     string
	AnswersIDS []string
}

type PollCreateRepository interface {
	CreatePoll(name string, answers []string) (*db.CreatePollResponse, error)
}

var ErrPollCreateNameEmpty = errors.New("poll name cannot be empty")
var ErrPollCreateAnswersLen = errors.New("poll needs at least 2 answers")
var ErrPollCreateAnswerEmpty = errors.New("poll answer name can't be empty")
var ErrPollCreateDB = errors.New("failed to create poll in db")

func (client *Client) CreatePoll(data *PollCreateRequest) (*PollCreateResponse, error) {
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

	dbPolls, err := client.PollCreateDB.CreatePoll(data.Name, data.Answers)
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

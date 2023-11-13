package poll

import (
	"errors"
	"strings"

	"github.com/vjerci/polls-app/pkg/domain/db"
)

var CreateMinAnswers = 2

type CreateRequest struct {
	Name    string
	Answers []string
}

type CreateResponse struct {
	PollID     string   `json:"poll_id"`
	AnswersIDS []string `json:"answer_ids"`
}

type CreateRepository interface {
	CreatePoll(name string, answers []string) (*db.PollCreateResponse, error)
}

type CreateModel struct {
	CreateDB CreateRepository
}

var ErrCreateNameEmpty = errors.New("poll name cannot be empty")
var ErrCreateAnswersLen = errors.New("poll needs at least 2 answers")
var ErrCreateAnswerEmpty = errors.New("poll answer name can't be empty")
var ErrCreateDB = errors.New("failed to create poll in db")

func (model *CreateModel) Create(data *CreateRequest) (*CreateResponse, error) {
	if data.Name == "" {
		return nil, ErrCreateNameEmpty
	}

	if len(data.Answers) < CreateMinAnswers {
		return nil, ErrCreateAnswersLen
	}

	for _, answer := range data.Answers {
		if answer == "" {
			return nil, ErrCreateAnswerEmpty
		}
	}

	pollName := strings.Trim(data.Name, " ")
	if !strings.HasSuffix(pollName, "?") {
		pollName += "?"
	}

	dbPolls, err := model.CreateDB.CreatePoll(pollName, data.Answers)
	if err != nil {
		return nil, errors.Join(ErrCreateDB, err)
	}

	resp := &CreateResponse{
		PollID:     dbPolls.ID,
		AnswersIDS: nil,
	}

	for _, answer := range dbPolls.Answers {
		resp.AnswersIDS = append(resp.AnswersIDS, answer.ID)
	}

	return resp, nil
}

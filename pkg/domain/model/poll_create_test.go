package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type MockPollCreateDB struct {
	InputPollName    string
	InputPollAnswers []string
	Response         *db.PollCreateResponse
	ResponseError    error
}

func (mock *MockPollCreateDB) CreatePoll(name string, answers []string) (*db.PollCreateResponse, error) {
	mock.InputPollName = name
	mock.InputPollAnswers = answers

	return mock.Response, mock.ResponseError
}

func TestPollCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *model.PollCreateRequest
		PollCreateDB  *MockPollCreateDB
	}{
		{
			ExpectedError: model.ErrPollCreateNameEmpty,
			Input: &model.PollCreateRequest{
				Name: "",
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: model.ErrPollCreateAnswersLen,
			Input: &model.PollCreateRequest{
				Name:    "pollName",
				Answers: nil,
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: model.ErrPollCreateAnswerEmpty,
			Input: &model.PollCreateRequest{
				Name:    "pollName",
				Answers: []string{"answerName", ""},
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: model.ErrPollCreateDB,
			Input: &model.PollCreateRequest{
				Name:    "pollName",
				Answers: []string{"answerName", "answerName1"},
			},
			PollCreateDB: &MockPollCreateDB{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		client := model.Client{
			PollCreateDB: test.PollCreateDB,
		}

		resp, err := client.CreatePoll(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestPollCreateSuccess(t *testing.T) {
	t.Parallel()

	pollCreateDBMock := &MockPollCreateDB{
		ResponseError: nil,
		Response: &db.PollCreateResponse{
			Name: "pollName",
			ID:   "pollID",
			Answers: []db.PollCreateAnswer{
				{
					Name: "answerName1",
					ID:   "answerID1",
				},
				{
					Name: "answerName2",
					ID:   "answerID2",
				},
			},
		},
	}

	client := model.Client{
		PollCreateDB: pollCreateDBMock,
	}

	input := &model.PollCreateRequest{
		Name:    "pollName",
		Answers: []string{"answer1", "answer2"},
	}

	resp, err := client.CreatePoll(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.Name, pollCreateDBMock.InputPollName,
		"expected input poll name to be passed to pollCreateDB")
	assert.EqualValues(t, input.Answers, pollCreateDBMock.InputPollAnswers,
		"expected input poll answers to be passed to pollCreateDB")

	assert.EqualValues(t, pollCreateDBMock.Response.ID, resp.PollID, "expected to return same poll id that db returned")

	for i, v := range pollCreateDBMock.Response.Answers {
		assert.EqualValues(t, resp.AnswersIDS[i], v.ID, "expected returned same answer id that db returned")
	}
}

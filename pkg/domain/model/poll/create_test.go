package poll_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
)

type MockCreateDB struct {
	InputName     string
	InputAnswers  []string
	Response      *db.PollCreateResponse
	ResponseError error
}

func (mock *MockCreateDB) CreatePoll(name string, answers []string) (*db.PollCreateResponse, error) {
	mock.InputName = name
	mock.InputAnswers = answers

	return mock.Response, mock.ResponseError
}

func TestCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *poll.CreateRequest
		PollCreateDB  *MockCreateDB
	}{
		{
			ExpectedError: poll.ErrCreateNameEmpty,
			Input: &poll.CreateRequest{
				Name: "",
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: poll.ErrCreateAnswersLen,
			Input: &poll.CreateRequest{
				Name:    "pollName",
				Answers: nil,
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: poll.ErrCreateAnswerEmpty,
			Input: &poll.CreateRequest{
				Name:    "pollName",
				Answers: []string{"answerName", ""},
			},
			PollCreateDB: nil,
		},
		{
			ExpectedError: poll.ErrCreateDB,
			Input: &poll.CreateRequest{
				Name:    "pollName",
				Answers: []string{"answerName", "answerName1"},
			},
			PollCreateDB: &MockCreateDB{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		CreateModel := poll.CreateModel{
			CreateDB: test.PollCreateDB,
		}

		resp, err := CreateModel.Create(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestCreateSuccess(t *testing.T) {
	t.Parallel()

	createDBMock := &MockCreateDB{
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

	createModel := poll.CreateModel{
		CreateDB: createDBMock,
	}

	input := &poll.CreateRequest{
		Name:    "pollName",
		Answers: []string{"answer1", "answer2"},
	}

	resp, err := createModel.Create(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.Name+"?", createDBMock.InputName,
		"expected input poll name to be passed to pollCreateDB")
	assert.EqualValues(t, input.Answers, createDBMock.InputAnswers,
		"expected input poll answers to be passed to pollCreateDB")

	assert.EqualValues(t, createDBMock.Response.ID, resp.PollID, "expected to return same poll id that db returned")

	for i, v := range createDBMock.Response.Answers {
		assert.EqualValues(t, resp.AnswersIDS[i], v.ID, "expected to return same answer ids that db returned")
	}
}

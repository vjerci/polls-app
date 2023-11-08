package poll_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
)

type MockDetailsDB struct {
	InputDetailsID     string
	ResponseDetailsErr error
	ResponseDetails    *db.PollDetailsResponse

	InputAnswersID     string
	ResponseAnswersErr error
	ResponseAnswers    []db.PollDetailsAnswer

	InputUserAnswerID     string
	InputUserAnswerUserID string
	ResponseUserAnswerErr error
	ResponseUserAnswer    string
}

func (mock *MockDetailsDB) GetPollDetails(pollID string) (*db.PollDetailsResponse, error) {
	mock.InputDetailsID = pollID

	return mock.ResponseDetails, mock.ResponseDetailsErr
}

func (mock *MockDetailsDB) GetPollDetailsAnswers(pollID string) ([]db.PollDetailsAnswer, error) {
	mock.InputAnswersID = pollID

	return mock.ResponseAnswers, mock.ResponseAnswersErr
}

func (mock *MockDetailsDB) GetUserAnswer(pollID string, userID string) (string, error) {
	mock.InputUserAnswerID = pollID
	mock.InputUserAnswerUserID = userID

	return mock.ResponseUserAnswer, mock.ResponseUserAnswerErr
}

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *poll.DetailsRequest
		PollDetailsDB *MockDetailsDB
	}{
		{
			ExpectedError: poll.ErrDetailsIDEmpty,
			Input: &poll.DetailsRequest{
				ID:     "",
				UserID: "userID",
			},
			PollDetailsDB: nil,
		},
		{
			ExpectedError: poll.ErrDetailsUserIDEmpty,
			Input: &poll.DetailsRequest{
				ID:     "pollID",
				UserID: "",
			},
			PollDetailsDB: nil,
		},
		{
			ExpectedError: poll.ErrDetailsNoPoll,
			Input: &poll.DetailsRequest{
				ID:     "pollID",
				UserID: "userID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails:    nil,
				ResponseDetailsErr: db.ErrPollDetailsNotFound,
			},
		},
		{
			ExpectedError: poll.ErrDetailsQueryInfo,
			Input: &poll.DetailsRequest{
				ID:     "pollID",
				UserID: "userID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails:    nil,
				ResponseDetailsErr: errors.New("test error"),
			},
		},
		{
			ExpectedError: poll.ErrDetailsAnswers,
			Input: &poll.DetailsRequest{
				ID:     "pollID",
				UserID: "userID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails:    &db.PollDetailsResponse{},
				ResponseDetailsErr: nil,

				ResponseAnswersErr: errors.New("test error"),
			},
		},
		{
			ExpectedError: poll.ErrDetailsUserAnswer,
			Input: &poll.DetailsRequest{
				ID:     "pollID",
				UserID: "userID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails:    &db.PollDetailsResponse{},
				ResponseDetailsErr: nil,

				ResponseAnswers:    []db.PollDetailsAnswer{},
				ResponseAnswersErr: nil,

				ResponseUserAnswer:    "",
				ResponseUserAnswerErr: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		pollDetailsModel := poll.DetailsModel{
			DetailsDB: test.PollDetailsDB,
		}

		resp, err := pollDetailsModel.Get(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestPollDetailsSuccess(t *testing.T) {
	t.Parallel()

	testsCases := []struct {
		Input         *poll.DetailsRequest
		PollDetailsDB *MockDetailsDB
	}{
		{
			Input: &poll.DetailsRequest{
				UserID: "userID",
				ID:     "pollID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails: &db.PollDetailsResponse{
					Name: "pollName",
					ID:   "pollID",
				},
				ResponseAnswers: []db.PollDetailsAnswer{
					{
						Name:  "answer1",
						ID:    "idAnswer1",
						Count: 2,
					},
					{
						Name:  "answer2",
						ID:    "idAnswer2",
						Count: 3,
					},
				},
				ResponseUserAnswer: "test",
			},
		},
		{
			Input: &poll.DetailsRequest{
				UserID: "userID",
				ID:     "pollID",
			},
			PollDetailsDB: &MockDetailsDB{
				ResponseDetails: &db.PollDetailsResponse{
					Name: "pollName",
					ID:   "pollID",
				},
				ResponseAnswers: []db.PollDetailsAnswer{
					{
						Name:  "answer1",
						ID:    "idAnswer1",
						Count: 2,
					},
					{
						Name:  "answer2",
						ID:    "idAnswer2",
						Count: 3,
					},
				},
				ResponseUserAnswerErr: db.ErrUserAnswerNotFound,
				ResponseUserAnswer:    "",
			},
		},
	}

	for _, test := range testsCases {
		pollDetailsModel := poll.DetailsModel{
			DetailsDB: test.PollDetailsDB,
		}

		resp, err := pollDetailsModel.Get(test.Input)

		if err != nil {
			t.Fatalf(`expected no err but got "%s" instead`, err)
		}

		assert.EqualValues(t,
			test.Input.ID,
			test.PollDetailsDB.InputDetailsID,
			"expected poll id to be passed to db when querying for details")
		assert.EqualValues(t,
			test.Input.ID,
			test.PollDetailsDB.InputAnswersID,
			"expected poll id to be passed to db when querying for answers")
		assert.EqualValues(t,
			test.Input.ID,
			test.PollDetailsDB.InputUserAnswerID,
			"expected poll id to be passed to db when querying for user answer")
		assert.EqualValues(t,
			test.Input.UserID,
			test.PollDetailsDB.InputUserAnswerUserID,
			"expected user id to be passed to db when querying for user answer")

		assert.EqualValues(t, test.PollDetailsDB.ResponseDetails.Name, resp.Name, "expected response poll name to match")
		assert.EqualValues(t, test.PollDetailsDB.ResponseDetails.ID, resp.ID, "expected response poll id to match")

		assert.EqualValues(t,
			test.PollDetailsDB.ResponseUserAnswer,
			resp.UserAnswer,
			"expected response poll user answer to match")

		for i, expectedAnswer := range test.PollDetailsDB.ResponseAnswers {
			assert.EqualValues(t, expectedAnswer.Count, resp.Answers[i].VotesCount, "expected answers votes count to match")
			assert.EqualValues(t, expectedAnswer.ID, resp.Answers[i].ID, "expected answers ids to match")
			assert.EqualValues(t, expectedAnswer.Name, resp.Answers[i].Name, "expected answers names to match")
		}
	}
}

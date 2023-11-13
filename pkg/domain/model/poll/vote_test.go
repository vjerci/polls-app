package poll_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
)

type MockVoteDB struct {
	InputAnswerExistsPollID   string
	InputAnswerExistsAnswerID string

	ResponseAnswerExistsErr error

	InputUpsertPollID   string
	InputUpsertAnswerID string
	InputUpsertUserID   string

	ResponseUpsert      bool
	ResponseUpsertError error
}

func (mock *MockVoteDB) ExistsPollAnswer(pollID, answerID string) error {
	mock.InputAnswerExistsPollID = pollID
	mock.InputAnswerExistsAnswerID = answerID

	return mock.ResponseAnswerExistsErr
}

func (mock *MockVoteDB) UpsertUserAnswer(pollID, answerID, userID string) (modified bool, err error) {
	mock.InputUpsertPollID = pollID
	mock.InputUpsertAnswerID = answerID
	mock.InputUpsertUserID = userID

	return mock.ResponseUpsert, mock.ResponseUpsertError
}

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *poll.VoteRequest
		PollVoteDB    *MockVoteDB
	}{
		{
			ExpectedError: poll.ErrVotePollIDEmpty,
			Input: &poll.VoteRequest{
				PollID: "",
			},
			PollVoteDB: nil,
		},
		{
			ExpectedError: poll.ErrVoteAnswerIDEmpty,
			Input: &poll.VoteRequest{
				PollID:   "pollID",
				AnswerID: "",
			},
			PollVoteDB: nil,
		},
		{
			ExpectedError: poll.ErrVoteUserIDEmpty,
			Input: &poll.VoteRequest{
				PollID:   "pollID",
				AnswerID: "answerID",
				UserID:   "",
			},
			PollVoteDB: nil,
		},
		{
			ExpectedError: poll.ErrVoteAnswerNotFound,
			Input: &poll.VoteRequest{
				PollID:   "pollID",
				AnswerID: "answerID",
				UserID:   "userID",
			},
			PollVoteDB: &MockVoteDB{
				ResponseAnswerExistsErr: errors.New("testError"),
			},
		},
		{
			ExpectedError: poll.ErrVoteUpsertFailed,
			Input: &poll.VoteRequest{
				PollID:   "pollID",
				AnswerID: "answerID",
				UserID:   "userID",
			},
			PollVoteDB: &MockVoteDB{
				ResponseAnswerExistsErr: nil,

				ResponseUpsert:      false,
				ResponseUpsertError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		pollVoteModel := poll.VoteModel{
			VoteDB: test.PollVoteDB,
		}

		resp, err := pollVoteModel.Do(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestPollVoteSuccess(t *testing.T) {
	t.Parallel()

	input := &poll.VoteRequest{
		PollID:   "testPollID",
		AnswerID: "testAnswerID",
		UserID:   "testUserID",
	}

	dbMock := MockVoteDB{
		ResponseAnswerExistsErr: nil,
		ResponseUpsert:          true,
		ResponseUpsertError:     nil,
	}

	pollVoteModel := poll.VoteModel{
		VoteDB: &dbMock,
	}

	resp, err := pollVoteModel.Do(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, dbMock.ResponseUpsert, resp.ModifiedAnswer, "expected to return modified from db upsert")

	assert.EqualValues(t,
		input.PollID,
		dbMock.InputAnswerExistsPollID,
		"expected to pass in poll id to dbRepository method ExistsPollAnswer")
	assert.EqualValues(t,
		input.AnswerID,
		dbMock.InputAnswerExistsAnswerID,
		"expected to pass in answer id to dbRepository method ExistsPollAnswer")

	assert.EqualValues(t,
		input.PollID,
		dbMock.InputUpsertPollID,
		"expected to pass in poll id to dbRepository method UpsertUserAnswer")
	assert.EqualValues(t,
		input.AnswerID,
		dbMock.InputUpsertAnswerID,
		"expected to pass in answer id to dbRepository method UpsertUserAnswer")
	assert.EqualValues(t,
		input.UserID,
		dbMock.InputUpsertUserID,
		"expected to pass in user id to dbRepository method UpsertUserAnswer")
}

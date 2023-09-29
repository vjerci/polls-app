package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type MockPollsListDB struct {
	InputPage     int
	Response      []db.PollListData
	ResponseError error
}

func (mock *MockPollsListDB) GetPollList(page int) ([]db.PollListData, error) {
	mock.InputPage = page

	return mock.Response, mock.ResponseError
}

type MockPollCountDB struct {
	InputPage int

	Response      bool
	ResponseError error
}

func (mock *MockPollCountDB) HasNextPage(page int) (bool, error) {
	mock.InputPage = page

	return mock.Response, mock.ResponseError
}

func TestPollsListErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *model.PollListRequest

		MockPollsListDB *MockPollsListDB
		MockPollCountDB *MockPollCountDB
	}{
		{
			ExpectedError: model.ErrPollListInvalidPage,
			Input: &model.PollListRequest{
				Page: -1,
			},
			MockPollsListDB: nil,
		},
		{
			ExpectedError: model.ErrPollListDB,
			Input: &model.PollListRequest{
				Page: 0,
			},
			MockPollsListDB: &MockPollsListDB{
				ResponseError: errors.New("test error"),
			},
		},
		{
			ExpectedError: model.ErrPollListNoPolls,
			Input: &model.PollListRequest{
				Page: 0,
			},
			MockPollsListDB: &MockPollsListDB{
				Response:      nil,
				ResponseError: nil,
			},
		},
		{
			ExpectedError: model.ErrPollListDBNextPage,
			Input: &model.PollListRequest{
				Page: 0,
			},
			MockPollsListDB: &MockPollsListDB{
				Response: []db.PollListData{
					{
						Name: "test",
						ID:   "testID",
					},
				},
				ResponseError: nil,
			},
			MockPollCountDB: &MockPollCountDB{
				Response:      false,
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		pollListModel := model.PollListModel{
			PollListDB:          test.MockPollsListDB,
			PollCountRepository: test.MockPollCountDB,
		}

		resp, err := pollListModel.Get(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestPollListSuccess(t *testing.T) {
	t.Parallel()

	pollListDBMock := &MockPollsListDB{
		ResponseError: nil,
		Response: []db.PollListData{
			{
				Name: "first name",
				ID:   "1",
			},
			{
				Name: "second name",
				ID:   "2",
			},
		},
	}

	pollCountDBMock := &MockPollCountDB{
		Response: true,
	}

	pollListModel := model.PollListModel{
		PollListDB:          pollListDBMock,
		PollCountRepository: pollCountDBMock,
	}

	input := &model.PollListRequest{
		Page: 1,
	}

	resp, err := pollListModel.Get(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.Page, pollListDBMock.InputPage, "expected input page to be passed to pollsListDB")
	assert.EqualValues(t, input.Page, pollCountDBMock.InputPage, "expected input page to be passed to pollCountDB")

	assert.EqualValues(t,
		pollCountDBMock.Response,
		resp.HasNext,
		"expected to pass hasNextPage got from database to response")

	for i, v := range pollListDBMock.Response {
		assert.EqualValues(t, resp.Polls[i].Name, v.Name, "expected returned Name to match response from pollsListDB")
		assert.EqualValues(t, resp.Polls[i].ID, v.ID, "expected returned ID to match response from pollsListDB")
	}
}

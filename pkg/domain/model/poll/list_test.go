package poll_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/db"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
)

type MockListDB struct {
	InputPage     int
	Response      []db.PollListData
	ResponseError error
}

func (mock *MockListDB) GetPollList(page int) ([]db.PollListData, error) {
	mock.InputPage = page

	return mock.Response, mock.ResponseError
}

type MockCountDB struct {
	InputPage int

	Response      bool
	ResponseError error
}

func (mock *MockCountDB) HasNextPage(page int) (bool, error) {
	mock.InputPage = page

	return mock.Response, mock.ResponseError
}

func TestPollsListErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *poll.ListRequest

		MockListDB  *MockListDB
		MockCountDB *MockCountDB
	}{
		{
			ExpectedError: poll.ErrListInvalidPage,
			Input: &poll.ListRequest{
				Page: -1,
			},
			MockListDB: nil,
		},
		{
			ExpectedError: poll.ErrListDB,
			Input: &poll.ListRequest{
				Page: 0,
			},
			MockListDB: &MockListDB{
				ResponseError: errors.New("test error"),
			},
		},
		{
			ExpectedError: poll.ErrListNoPolls,
			Input: &poll.ListRequest{
				Page: 0,
			},
			MockListDB: &MockListDB{
				Response:      nil,
				ResponseError: nil,
			},
		},
		{
			ExpectedError: poll.ErrListDBNextPage,
			Input: &poll.ListRequest{
				Page: 0,
			},
			MockListDB: &MockListDB{
				Response: []db.PollListData{
					{
						Name: "test",
						ID:   "testID",
					},
				},
				ResponseError: nil,
			},
			MockCountDB: &MockCountDB{
				Response:      false,
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		listModel := poll.ListModel{
			ListDB:  test.MockListDB,
			CountDB: test.MockCountDB,
		}

		resp, err := listModel.Get(test.Input)

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

	listDBMock := &MockListDB{
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

	countDBMock := &MockCountDB{
		Response: true,
	}

	pollListModel := poll.ListModel{
		ListDB:  listDBMock,
		CountDB: countDBMock,
	}

	input := &poll.ListRequest{
		Page: 1,
	}

	resp, err := pollListModel.Get(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.Page, listDBMock.InputPage, "expected input page to be passed to pollsListDB")
	assert.EqualValues(t, input.Page, countDBMock.InputPage, "expected input page to be passed to pollCountDB")

	assert.EqualValues(t,
		countDBMock.Response,
		resp.HasNext,
		"expected to pass hasNextPage got from database to response")

	for i, v := range listDBMock.Response {
		assert.EqualValues(t, resp.Polls[i].Name, v.Name, "expected returned Name to match response from pollsListDB")
		assert.EqualValues(t, resp.Polls[i].ID, v.ID, "expected returned ID to match response from pollsListDB")
	}
}

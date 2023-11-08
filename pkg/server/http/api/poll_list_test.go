package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type MockPollsListModel struct {
	InputData     *poll.ListRequest
	ResponseData  *poll.ListResponse
	ResponseError error
}

func (mock *MockPollsListModel) Get(input *poll.ListRequest) (*poll.ListResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollListErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		URL           string
		Model         *MockPollsListModel
	}{
		{
			ExpectedError: schema.ErrPollListPageNotInt,
			URL:           "http://localhost/poll_list?page=dsadsa",
		},
		{
			ExpectedError: schema.ErrPollListModel,
			URL:           "http://localhost/poll_list?page=1",
			Model: &MockPollsListModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, test.URL, nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		echoContext := e.NewContext(req, rec)

		apiClient := api.New(&api.Models{
			PollList: test.Model,
		}, &api.SchemaMap{
			PollList: &schema.PollListSchemaMap{},
		})

		err := apiClient.PollList(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollListSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"page":4}`
	pollListModelMock := &MockPollsListModel{
		ResponseData: &poll.ListResponse{
			Polls: []poll.GeneralInfo{
				{
					Name: "Do you want a lift?",
					ID:   "1",
				},
				{
					Name: "do you want a lightning?",
					ID:   "2",
				},
			},
			HasNext: true,
		},
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":{"polls":` +
		`[{"name":"Do you want a lift?","id":"1"},{"name":"do you want a lightning?","id":"2"}]` +
		`,"has_next":true` + `}}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/polls_list", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	echoContext := e.NewContext(req, rec)

	apiClient := api.New(&api.Models{
		PollList: pollListModelMock,
	}, &api.SchemaMap{
		PollList: &schema.PollListSchemaMap{},
	})

	err := apiClient.PollList(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

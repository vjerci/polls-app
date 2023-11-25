package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type MockPollDetailsModel struct {
	InputData     *poll.DetailsRequest
	ResponseData  *poll.DetailsResponse
	ResponseError error
}

func (mock *MockPollDetailsModel) Get(input *poll.DetailsRequest) (*poll.DetailsResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         func(echoContext echo.Context) echo.Context
		Model         *MockPollDetailsModel
	}{
		{
			ExpectedError: api.ErrUserIDIsNotString,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", struct{}{})

				return echoContext
			},
		},
		{
			ExpectedError: schema.ErrPollDetailsModel,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", "testUserID")

				return echoContext
			},
			Model: &MockPollDetailsModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/poll_list", strings.NewReader(""))
		rec := httptest.NewRecorder()

		e := echo.New()
		echoContext := e.NewContext(req, rec)

		echoContext = test.Input(echoContext)

		apiClient := api.New(&api.Models{
			PollDetails: test.Model,
		}, &api.SchemaMap{
			PollDetails: &schema.PollDetailsSchemaMap{},
		})

		err := apiClient.PollDetails(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollDetailsSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"page":4}`
	pollLDetailsModelMock := &MockPollDetailsModel{
		ResponseData: &poll.DetailsResponse{
			ID:         "pollID",
			Name:       "pollName",
			UserAnswer: "answerID",
			Answers: []poll.DetailsAnswer{
				{
					ID:         "answerID",
					Name:       "answerName",
					VotesCount: 2,
				},
			},
		},
		ResponseError: nil,
	}

	req := httptest.NewRequest(echo.POST, "http://localhost/polls_list", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	echoContext := e.NewContext(req, rec)

	echoContext.Set("userID", "testUserID")

	apiClient := api.New(&api.Models{
		PollDetails: pollLDetailsModelMock,
	}, &api.SchemaMap{
		PollDetails: &schema.PollDetailsSchemaMap{},
	})

	err := apiClient.PollDetails(echoContext)

	expectedResponse := `{"success":true,"data":{"id":"pollID","name":"pollName","user_vote":"answerID"` +
		`,"answers":[{"id":"answerID","name":"answerName","votes_count":2}` +
		`]}}` + "\n"

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

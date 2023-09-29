package api_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type MockPollVoteModel struct {
	InputData *model.PollVoteRequest

	ResponseData  *model.PollVoteResponse
	ResponseError error
}

func (mock *MockPollVoteModel) Do(input *model.PollVoteRequest) (*model.PollVoteResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         func(echoContext echo.Context) echo.Context
		Model         *MockPollVoteModel
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
			ExpectedError: schema.ErrPollVoteJSONDecode,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Request().Body = io.NopCloser(strings.NewReader("[{]}"))

				echoContext.Set("userID", "testUserID")

				return echoContext
			},
		},
		{
			ExpectedError: schema.ErrPollVoteModel,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", "testUserID")

				return echoContext
			},
			Model: &MockPollVoteModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/poll/vote", strings.NewReader(`{"answer_id":"testID"}`))
		rec := httptest.NewRecorder()

		e := echo.New()
		echoContext := e.NewContext(req, rec)

		echoContext = test.Input(echoContext)

		apiClient := api.New(&api.Models{
			PollVote: test.Model,
		}, &api.SchemaMap{
			PollVote: &schema.PollVoteSchemaMap{},
		})

		err := apiClient.PollVote(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollVoteSuccessful(t *testing.T) {
	t.Parallel()

	pollID := "testPollID"
	answerID := "testAnswerID"
	userID := "testUserID"

	modelMock := &MockPollVoteModel{
		ResponseData: &model.PollVoteResponse{
			ModifiedAnswer: true,
		},
	}

	req := httptest.NewRequest(echo.GET, "http://localhost/poll/vote", strings.NewReader(`{"answer_id":"`+answerID+`"}`))
	rec := httptest.NewRecorder()

	e := echo.New()
	echoContext := e.NewContext(req, rec)
	echoContext.SetParamNames("id")
	echoContext.SetParamValues(pollID)

	echoContext.Set("userID", userID)

	apiClient := api.New(&api.Models{
		PollVote: modelMock,
	}, &api.SchemaMap{
		PollVote: &schema.PollVoteSchemaMap{},
	})

	err := apiClient.PollVote(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	expectedResponse := `{"success":true,"data":{"modified_answer":true}}` + "\n"

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")

	assert.EqualValues(t, pollID, modelMock.InputData.PollID, "expected user id to be passed to model")
	assert.EqualValues(t, answerID, modelMock.InputData.AnswerID, "expected answer id to be passed to model")
	assert.EqualValues(t, userID, modelMock.InputData.UserID, "expected user id to be passed to model")
}

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

func (mock *MockPollVoteModel) PollVote(input *model.PollVoteRequest) (*model.PollVoteResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         func(echoContext echo.Context) echo.Context
		Model         *MockPollVoteModel
		ErrorMap      api.PollVoteSchemaMap
	}{
		{
			ExpectedError: api.ErrUserIDIsNotString,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", struct{}{})

				return echoContext
			},
			ErrorMap: schema.NewSchemaMap(),
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
			ErrorMap: schema.NewSchemaMap(),
		},
		{
			ExpectedError: schema.ErrPollVoteModel,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", "testUserID")

				return echoContext
			},
			ErrorMap: schema.NewSchemaMap(),
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

		factory := api.New()

		err := factory.PollVote(test.Model, test.ErrorMap)(echoContext)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestPollVoteSuccessful(t *testing.T) {
	t.Parallel()

	pollID := "testPollID"
	answerID := "testAnswerID"
	userID := "testUserID"

	model := &MockPollVoteModel{
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

	factory := api.New()

	err := factory.PollVote(model, schema.NewSchemaMap())(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	expectedResponse := `{"success":true,"data":{"modified_answer":true}}` + "\n"

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")

	assert.EqualValues(t, pollID, model.InputData.PollID, "expected user id to be passed to model")
	assert.EqualValues(t, answerID, model.InputData.AnswerID, "expected answer id to be passed to model")
	assert.EqualValues(t, userID, model.InputData.UserID, "expected user id to be passed to model")
}

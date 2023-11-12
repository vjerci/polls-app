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
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type MockPollCreateModel struct {
	InputData     *poll.CreateRequest
	ResponseData  *poll.CreateResponse
	ResponseError error
}

func (mock *MockPollCreateModel) Create(input *poll.CreateRequest) (*poll.CreateResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *strings.Reader
		Model         *MockPollCreateModel
	}{
		{
			ExpectedError: schema.ErrPollCreateJSONDecode,
			Input:         strings.NewReader("[{]}"),
		},
		{
			ExpectedError: schema.ErrPollCreateModel,
			Input:         strings.NewReader(`{"name":"testName","answers":["answer1","answer2"]}`),
			Model: &MockPollCreateModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.PUT, "http://localhost/poll", test.Input)
		rec := httptest.NewRecorder()

		e := echo.New()
		echoContext := e.NewContext(req, rec)

		apiClient := api.New(&api.Models{
			PollCreate: test.Model,
		}, &api.SchemaMap{
			PollCreate: &schema.PollCreateSchemaMap{},
		})

		err := apiClient.PollCreate(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollCreateSuccessful(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`{"name":"testName","answers":["answer1","answer2"]}`)
	pollCreateModelMock := &MockPollCreateModel{
		ResponseData: &poll.CreateResponse{
			PollID:     "pollID",
			AnswersIDS: []string{"answerID1", "answerID2"},
		},
		ResponseError: nil,
	}

	req := httptest.NewRequest(echo.PUT, "http://localhost/poll", input)
	rec := httptest.NewRecorder()

	e := echo.New()
	echoContext := e.NewContext(req, rec)

	apiClient := api.New(&api.Models{
		PollCreate: pollCreateModelMock,
	}, &api.SchemaMap{
		PollCreate: &schema.PollCreateSchemaMap{},
	})

	err := apiClient.PollCreate(echoContext)

	expectedResponse := `{"success":true,"data":{"id":"pollID","answers_ids":["answerID1","answerID2"]}}` + "\n"

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

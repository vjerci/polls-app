package api_test

import (
	"errors"
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

type MockPollDetailsModel struct {
	InputData     *model.PollDetailsRequest
	ResponseData  *model.PollDetailsResponse
	ResponseError error
}

func (mock *MockPollDetailsModel) GetPollDetails(input *model.PollDetailsRequest) (*model.PollDetailsResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         func(echoContext echo.Context) echo.Context
		Model         *MockPollDetailsModel
		ErrorMap      api.PollDetailsSchemaMap
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
			ExpectedError: schema.ErrPollDetailsModel,
			Input: func(echoContext echo.Context) echo.Context {
				echoContext.SetParamNames("id")
				echoContext.SetParamValues("testPollID")

				echoContext.Set("userID", "testUserID")

				return echoContext
			},
			ErrorMap: schema.NewSchemaMap(),
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

		factory := api.New()

		err := factory.PollDetails(test.Model, test.ErrorMap)(echoContext)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestPollDetailsSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"page":4}`
	pollListModelMock := &MockPollsListModel{
		ResponseData: &model.PollListResponse{
			Polls: []model.GeneralPollInfo{
				{
					Name: "Do you want a lift?",
					ID:   "1",
				},
				{
					Name: "do you want a lightning?",
					ID:   "2",
				},
			},
		},
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":{"polls":[` +
		`{"name":"Do you want a lift?","id":"1"},{"name":"do you want a lightning?","id":"2"}` +
		`]}}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/polls_list", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	factory := api.New()

	err := factory.PollList(pollListModelMock, schema.NewSchemaMap())(c)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

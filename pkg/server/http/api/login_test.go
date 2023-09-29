package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type MockLoginModel struct {
	InputData     *model.LoginRequest
	ResponseData  *model.LoginResponse
	ResponseError error
}

func (mock *MockLoginModel) Do(input *model.LoginRequest) (*model.LoginResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         string
		Model         *MockLoginModel
	}{
		{
			ExpectedError: schema.ErrLoginJSONDecode,
			Input:         "[{]}",
		},
		{
			ExpectedError: schema.ErrLoginModel,
			Input:         `{}`,
			Model: &MockLoginModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/login", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		echoInstance := echo.New()
		echoContext := echoInstance.NewContext(req, rec)

		apiClient := api.New(&api.Models{
			Login: test.Model,
		}, &api.SchemaMap{
			Login: &schema.LoginSchemaMap{},
		})

		err := apiClient.Login(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestLoginSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"user_id":"userID"}`
	loginModelMock := &MockLoginModel{
		ResponseData: &model.LoginResponse{
			Token: "testToken",
			Name:  "Jhon",
		},
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":{"token":"testToken","name":"Jhon"}}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/login", strings.NewReader(input))
	rec := httptest.NewRecorder()

	echoInstance := echo.New()
	echoContext := echoInstance.NewContext(req, rec)

	apiClient := api.New(&api.Models{
		Login: loginModelMock,
	}, &api.SchemaMap{
		Login: &schema.LoginSchemaMap{},
	})

	err := apiClient.Login(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

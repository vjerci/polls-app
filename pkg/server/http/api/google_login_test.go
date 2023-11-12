package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type MockGoogleLoginModel struct {
	InputData     *auth.GoogleLoginRequest
	ResponseData  *auth.GoogleLoginResponse
	ResponseError error
}

func (mock *MockGoogleLoginModel) Do(input *auth.GoogleLoginRequest) (*auth.GoogleLoginResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestGoogleLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         string
		Model         *MockGoogleLoginModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginJSONDecode,
			Input:         "[{]}",
		},
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input:         `{}`,
			Model: &MockGoogleLoginModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.POST, "http://localhost/google/login", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		echoInstance := echo.New()
		echoContext := echoInstance.NewContext(req, rec)

		apiClient := api.New(&api.Models{
			GoogleLogin: test.Model,
		}, &api.SchemaMap{
			GoogleLogin: &schema.GoogleLoginSchemaMap{},
		})

		err := apiClient.GoogleLogin(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestGoogleLoginSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"user_id":"userID"}`
	googleLoginModelMock := &MockGoogleLoginModel{
		ResponseData: &auth.GoogleLoginResponse{
			Token: "testToken",
			Name:  "Jhon",
		},
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":{"token":"testToken","name":"Jhon"}}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/gooogle/login", strings.NewReader(input))
	rec := httptest.NewRecorder()

	echoInstance := echo.New()
	echoContext := echoInstance.NewContext(req, rec)

	apiClient := api.New(&api.Models{
		GoogleLogin: googleLoginModelMock,
	}, &api.SchemaMap{
		GoogleLogin: &schema.GoogleLoginSchemaMap{},
	})

	err := apiClient.GoogleLogin(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	modelauth "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type MockRegisterModel struct {
	InputData     *modelauth.RegisterRequest
	ResponseData  auth.AccessToken
	ResponseError error
}

func (mock *MockRegisterModel) Do(input *modelauth.RegisterRequest) (accessToken auth.AccessToken, err error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         string
		Model         *MockRegisterModel
	}{
		{
			ExpectedError: schema.ErrRegisterJSONDecode,
			Input:         "[{]}",
		},
		{
			ExpectedError: schema.ErrRegisterModel,
			Input:         `{}`,
			Model: &MockRegisterModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/register", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		e := echo.New()
		echoContext := e.NewContext(req, rec)

		apiClient := api.New(&api.Models{
			Register: test.Model,
		}, &api.SchemaMap{
			Register: &schema.RegisterSchemaMap{},
		})

		err := apiClient.Register(echoContext)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestRegisterSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"user_id":"userID","name":"name"}`
	registerModelMock := &MockRegisterModel{
		ResponseData:  auth.AccessToken("token"),
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":"token"}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/login", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	echoContext := e.NewContext(req, rec)

	apiClient := api.New(&api.Models{
		Register: registerModelMock,
	}, &api.SchemaMap{
		Register: &schema.RegisterSchemaMap{},
	})

	err := apiClient.Register(echoContext)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

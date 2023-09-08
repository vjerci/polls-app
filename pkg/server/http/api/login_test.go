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

type MockLoginModel struct {
	InputData     *model.LoginRequest
	ResponseData  *model.LoginResponse
	ResponseError error
}

func (mock *MockLoginModel) Login(input *model.LoginRequest) (*model.LoginResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         string
		Model         *MockLoginModel
		ErrorMap      api.LoginSchemaMap
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
			ErrorMap: schema.NewSchemaMap(),
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/login", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		factory := api.New()

		err := factory.Login(test.Model, test.ErrorMap)(c)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
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

	e := echo.New()
	c := e.NewContext(req, rec)

	factory := api.New()

	err := factory.Login(loginModelMock, schema.NewSchemaMap())(c)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

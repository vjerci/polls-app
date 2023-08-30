package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type MockLoginModel struct {
	ResponseToken login.AccessToken
	ResponseName  string
	ResponseError error
	LoginData     model.LoginData
}

func (mock *MockLoginModel) Login(data model.LoginData) (login.AccessToken, string, error) {
	mock.LoginData = data
	return mock.ResponseToken, mock.ResponseName, mock.ResponseError
}

func TestLoginErrorHandling(t *testing.T) {
	testCases := []struct {
		ExpectedError error
		Input         string
		Model         *MockLoginModel
	}{
		{
			ExpectedError: api.ErrLoginJsonDecode,
			Input:         "[{]}",
		},
		{
			ExpectedError: api.ErrLoginUserDoesNotExist,
			Input:         `{}`,
			Model: &MockLoginModel{
				ResponseToken: "",
				ResponseError: model.ErrLoginUserNotFound,
			},
		},
		{
			ExpectedError: api.ErrLoginUserIDNotSet,
			Input:         `{}`,
			Model: &MockLoginModel{
				ResponseToken: "",
				ResponseError: model.ErrLoginUserIDNotSet,
			},
		},
		{
			ExpectedError: api.ErrorLoginModel,
			Input:         `{}`,
			Model: &MockLoginModel{
				ResponseToken: "",
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/register", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		factory := api.New()

		err := factory.Login(test.Model)(c)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestLoginSuccessful(t *testing.T) {
	input := `{"user_id":"userID"}`
	model := &MockLoginModel{
		ResponseToken: login.AccessToken("testToken"),
		ResponseName:  "Jhon",
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":{"name":"Jhon","access_token":"testToken"}}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/login", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	factory := api.New()

	err := factory.Login(model)(c)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

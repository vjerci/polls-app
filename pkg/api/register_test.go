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

type MockRegisterModel struct {
	ResponseToken    login.AccessToken
	ResponseError    error
	RegistrationData model.RegistrationData
}

func (mock *MockRegisterModel) Register(data model.RegistrationData) (login.AccessToken, error) {
	mock.RegistrationData = data
	return mock.ResponseToken, mock.ResponseError
}

func TestRegisterErrorHandling(t *testing.T) {
	testCases := []struct {
		ExpectedError error
		Input         string
		Model         *MockRegisterModel
	}{
		{
			ExpectedError: api.ErrRegisterJsonDecode,
			Input:         "[{]}",
		},
		{
			ExpectedError: api.ErrRegisterUserIDNotSet,
			Input:         `{}`,
			Model: &MockRegisterModel{
				ResponseToken: "",
				ResponseError: model.ErrRegisterUserIDNotSet,
			},
		},
		{
			ExpectedError: api.ErrRegisterNameNotSet,
			Input:         `{}`,
			Model: &MockRegisterModel{
				ResponseToken: "",
				ResponseError: model.ErrRegisterNameNotSet,
			},
		},
		{
			ExpectedError: api.ErrRegisterUserDuplicate,
			Input:         `{}`,
			Model: &MockRegisterModel{
				ResponseToken: "",
				ResponseError: model.ErrRegisterDuplicate,
			},
		},
		{
			ExpectedError: api.ErrRegisterModel,
			Input:         `{}`,
			Model: &MockRegisterModel{
				ResponseToken: "",
				ResponseError: errors.New("generic error from db"),
			},
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/register", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		factory := api.New()

		err := factory.Register(test.Model)(c)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestRegisterSuccessful(t *testing.T) {
	input := `{"user_id":"userID","group_id":"groupID","name":"name"}`
	token := login.AccessToken("testToken")
	model := &MockRegisterModel{
		ResponseToken: token,
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":"testToken"}` + "\n"

	req := httptest.NewRequest(echo.PUT, "http://localhost/register", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	factory := api.New()

	err := factory.Register(model)(c)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

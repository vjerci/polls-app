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
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type MockRegisterModel struct {
	InputData     *model.RegisterRequest
	ResponseData  login.AccessToken
	ResponseError error
}

func (mock *MockRegisterModel) Register(input *model.RegisterRequest) (accessToken login.AccessToken, err error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         string
		Model         *MockRegisterModel
		ErrorMap      api.RegisterSchemaMap
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
			ErrorMap: schema.NewSchemaMap(),
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest(echo.GET, "http://localhost/register", strings.NewReader(test.Input))
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		factory := api.New()

		err := factory.Register(test.Model, test.ErrorMap)(c)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestRegisterSuccessful(t *testing.T) {
	t.Parallel()

	input := `{"user_id":"userID","name":"name"}`
	registerModelMock := &MockRegisterModel{
		ResponseData:  login.AccessToken("token"),
		ResponseError: nil,
	}
	expectedResponse := `{"success":true,"data":"token"}` + "\n"

	req := httptest.NewRequest(echo.POST, "http://localhost/login", strings.NewReader(input))
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	factory := api.New()

	err := factory.Register(registerModelMock, schema.NewSchemaMap())(c)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
}

package api_test

// import (
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	echo "github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
// 	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
// 	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/connector"
// 	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
// )

// type MockRegisterConnector struct {
// 	ResponseData  string
// 	ResponseError error
// 	InputData     model.RegistrationData
// 	RegisterModel connector.RegisterModel
// }

// func (mock *MockRegisterConnector) Register(
// 	registerModel connector.RegisterModel,
// 	input model.RegistrationData,
// ) (string, error) {
// 	mock.InputData = input
// 	mock.RegisterModel = registerModel

// 	return mock.ResponseData, mock.ResponseError
// }

// type MockRegisterModel struct{}

// func (mock *MockRegisterModel) Register(_ model.RegistrationData) (login.AccessToken, error) {
// 	return login.AccessToken(""), nil
// }

// func TestRegisterErrors(t *testing.T) {
// 	t.Parallel()

// 	testError := errors.New("test error")
// 	testCases := []struct {
// 		ExpectedError error
// 		Input         string
// 		Connector     *MockRegisterConnector
// 	}{
// 		{
// 			ExpectedError: api.ErrRegisterJSONDecode,
// 			Input:         "[{]}",
// 		},
// 		{
// 			ExpectedError: testError,
// 			Input:         `{}`,
// 			Connector: &MockRegisterConnector{
// 				ResponseData:  "",
// 				ResponseError: testError,
// 			},
// 		},
// 	}

// 	for _, test := range testCases {
// 		req := httptest.NewRequest(echo.GET, "http://localhost/register", strings.NewReader(test.Input))
// 		rec := httptest.NewRecorder()

// 		e := echo.New()
// 		c := e.NewContext(req, rec)

// 		factory := api.New()

// 		err := factory.Register(&MockRegisterModel{}, test.Connector)(c)

// 		if !errors.Is(err, test.ExpectedError) {
// 			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
// 		}
// 	}
// }

// func TestRegisterSuccessful(t *testing.T) {
// 	t.Parallel()

// 	input := `{"user_id":"userID","group_id":"groupID","name":"name"}`
// 	connector := &MockRegisterConnector{
// 		ResponseData:  "testToken",
// 		ResponseError: nil,
// 	}
// 	expectedResponse := `{"success":true,"data":"testToken"}` + "\n"

// 	req := httptest.NewRequest(echo.PUT, "http://localhost/register", strings.NewReader(input))
// 	rec := httptest.NewRecorder()

// 	e := echo.New()
// 	c := e.NewContext(req, rec)

// 	factory := api.New()

// 	err := factory.Register(&MockRegisterModel{}, connector)(c)

// 	if err != nil {
// 		t.Fatalf(`expected no err but got "%s" instead`, err)
// 	}

// 	assert.Equal(t, http.StatusOK, rec.Code, "response code doesn't match")
// 	assert.EqualValues(t, expectedResponse, rec.Body.String(), "response body doesn't match")
// }

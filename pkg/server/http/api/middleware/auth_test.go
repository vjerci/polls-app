package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	authmodel "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api/middleware"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type MockLoginRepo struct {
	Input auth.AccessToken

	Response    string
	ResponseErr error
}

func (mock *MockLoginRepo) Decode(input auth.AccessToken) (userID string, err error) {
	mock.Input = input

	return mock.Response, mock.ResponseErr
}

type MockUserRepo struct {
	Input string

	Response    string
	ResponseErr error
}

func (mock *MockUserRepo) GetUser(userID string) (name string, err error) {
	mock.Input = userID

	return mock.Response, mock.ResponseErr
}

func TestAuth(t *testing.T) {
	t.Parallel()

	mockError := errors.New("test")

	testCases := []struct {
		Name string

		Headers  http.Header
		AuthMock middleware.AuthRepo
		UserMock middleware.UserRepo

		ExpectedError error
		ExpectedField string
	}{
		{
			Name: "Authentication header missing",

			Headers: nil,

			ExpectedError: schema.ErrAuthMissing,
		},
		{
			Name: "Authentication header invalid",

			Headers: http.Header{
				"Authorization": []string{"test"},
			},
			ExpectedError: schema.ErrAuthInvalid,
		},
		{
			Name: "Auth repo error",

			Headers: http.Header{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				ResponseErr: errors.New("test error"),
			},

			ExpectedError: schema.ErrAuthInvalid,
		},
		{
			Name: "User repo - user not found failure",

			Headers: http.Header{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock: &MockUserRepo{
				ResponseErr: authmodel.ErrGetUserUserNotFound,
			},

			ExpectedError: schema.ErrAuthUserDoesNotExist,
		},
		{
			Name: "User repo - other types of failure",

			Headers: http.Header{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock: &MockUserRepo{
				ResponseErr: mockError,
			},

			ExpectedError: mockError,
		},
		{
			Name: "Passing test",

			Headers: http.Header{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock: &MockUserRepo{},

			ExpectedField: "testUserID1",
			ExpectedError: nil,
		},
	}

	for _, test := range testCases {
		testCase := test

		t.Run(t.Name()+"-"+testCase.Name, func(t *testing.T) {
			t.Parallel()

			echoServer := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			res := httptest.NewRecorder()
			req.Header = testCase.Headers
			echoContext := echoServer.NewContext(req, res)

			client := middleware.Client{
				AuthRepo: testCase.AuthMock,
				UserRepo: testCase.UserMock,
			}

			handler := client.WithAuth(echo.HandlerFunc(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}))
			err := handler(echoContext)

			assert.EqualValues(t, testCase.ExpectedError, err, "expected errors to match")

			if testCase.ExpectedField != "" {
				assert.EqualValues(t, testCase.ExpectedField, echoContext.Get("userID"), "expected userID to be set")
			}
		})
	}
}

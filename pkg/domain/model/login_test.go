package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
)

type MockLoginDB struct {
	UserID              string
	ResponseAccessToken login.AccessToken
	ResponseError       error
}

func (mock *MockLoginDB) CreateToken(userID string) (login.AccessToken, error) {
	mock.UserID = userID

	return mock.ResponseAccessToken, mock.ResponseError
}

type MockUserDB struct {
	UserID        string
	ResponseName  string
	ResponseError error
}

func (mock *MockUserDB) GetUser(userID string) (string, error) {
	mock.UserID = userID

	return mock.ResponseName, mock.ResponseError
}

func TestLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		Input         *model.LoginRequest
		UserDBMock    *MockUserDB
		LoginDBMock   *MockLoginDB
	}{
		{
			ExpectedError: model.ErrLoginUserIDNotSet,
			Input: &model.LoginRequest{
				UserID: "",
			},
			UserDBMock:  &MockUserDB{},
			LoginDBMock: &MockLoginDB{},
		},
		{
			ExpectedError: model.ErrLoginUserNotFound,
			Input: &model.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "",
				ResponseError: db.ErrGetUserNoRows,
			},
			LoginDBMock: &MockLoginDB{},
		},
		{
			ExpectedError: model.ErrLoginUserGet,
			Input: &model.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "",
				ResponseError: errors.New("test error"),
			},
			LoginDBMock: &MockLoginDB{},
		},
		{
			ExpectedError: model.ErrLoginCreateToken,
			Input: &model.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "userName",
				ResponseError: nil,
			},
			LoginDBMock: &MockLoginDB{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		client := model.Client{
			LoginDB: test.LoginDBMock,
			UserDB:  test.UserDBMock,
		}

		resp, err := client.Login(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got  %v instead", resp)
		}
	}
}

func TestLoginSuccess(t *testing.T) {
	t.Parallel()

	userDBMock := &MockUserDB{
		ResponseError: nil,
		ResponseName:  "Jhon",
	}

	loginDBMock := &MockLoginDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	client := model.Client{
		UserDB:  userDBMock,
		LoginDB: loginDBMock,
	}

	input := &model.LoginRequest{
		UserID: "userID",
	}

	resp, err := client.Login(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.UserID, userDBMock.UserID, "expected input user_id to be passed to userDB")
	assert.EqualValues(t, loginDBMock.UserID, input.UserID, "expected input's user_id to be passed to loginDB")

	assert.EqualValues(t, userDBMock.ResponseName, resp.Name, "expected returned name to match response from userDB")
	assert.EqualValues(t, loginDBMock.ResponseAccessToken, resp.Token, "expected token to match response from loginDB")
}

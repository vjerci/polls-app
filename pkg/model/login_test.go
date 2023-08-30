package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
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

func TestLoginErrorHandling(t *testing.T) {
	testCases := []struct {
		ExpectedError error
		Input         model.LoginData
		UserDBMock    *MockUserDB
		LoginDBMock   *MockLoginDB
	}{
		{
			ExpectedError: model.ErrLoginUserIDNotSet,
			Input: model.LoginData{
				UserID: "",
			},
			UserDBMock:  &MockUserDB{},
			LoginDBMock: &MockLoginDB{},
		},
		{
			ExpectedError: model.ErrLoginUserNotFound,
			Input: model.LoginData{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "",
				ResponseError: db.ErrGetUserNoRows,
			},
			LoginDBMock: &MockLoginDB{},
		},
		{
			ExpectedError: model.ErrLoginUserGetUser,
			Input: model.LoginData{
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
			Input: model.LoginData{
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

		token, name, err := client.Login(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		assert.EqualValues(t, "", token, "expected token to be empty")
		assert.EqualValues(t, "", name, "expected name to be empty")
	}
}

func TestLoginSuccess(t *testing.T) {
	userDBMock := &MockUserDB{
		ResponseError: nil,
		ResponseName:  "Jhon",
	}

	loginDbMock := &MockLoginDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	client := model.Client{
		UserDB:  userDBMock,
		LoginDB: loginDbMock,
	}

	input := model.LoginData{
		UserID: "userID",
	}

	token, name, err := client.Login(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.UserID, userDBMock.UserID, "expected input user_id to be passed to userDB")
	assert.EqualValues(t, loginDbMock.UserID, input.UserID, "expected input's user_id to be passed to loginDB")

	assert.EqualValues(t, userDBMock.ResponseName, name, "expected returned name to match response from userDB")
	assert.EqualValues(t, loginDbMock.ResponseAccessToken, token, "expected token to match response from loginDB")
}

package auth_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/db"
	modelauth "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
)

type MockAuthDB struct {
	UserID              string
	ResponseAccessToken auth.AccessToken
	ResponseError       error
}

func (mock *MockAuthDB) CreateToken(userID string) (auth.AccessToken, error) {
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
		Input         *modelauth.LoginRequest
		UserDBMock    *MockUserDB
		AuthDBMock    *MockAuthDB
	}{
		{
			ExpectedError: modelauth.ErrLoginUserIDNotSet,
			Input: &modelauth.LoginRequest{
				UserID: "",
			},
			UserDBMock: &MockUserDB{},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrLoginUserNotFound,
			Input: &modelauth.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "",
				ResponseError: db.ErrGetUserNoRows,
			},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrLoginUserDB,
			Input: &modelauth.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "",
				ResponseError: errors.New("test error"),
			},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrLoginCreateToken,
			Input: &modelauth.LoginRequest{
				UserID: "userID",
			},
			UserDBMock: &MockUserDB{
				ResponseName:  "userName",
				ResponseError: nil,
			},
			AuthDBMock: &MockAuthDB{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		loginModel := modelauth.LoginModel{
			AuthDB: test.AuthDBMock,
			UserDB: test.UserDBMock,
		}

		resp, err := loginModel.Do(test.Input)

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

	AuthDBMock := &MockAuthDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	loginModel := modelauth.LoginModel{
		UserDB: userDBMock,
		AuthDB: AuthDBMock,
	}

	input := &modelauth.LoginRequest{
		UserID: "userID",
	}

	resp, err := loginModel.Do(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.UserID, userDBMock.UserID, "expected input user_id to be passed to userDB")
	assert.EqualValues(t, AuthDBMock.UserID, input.UserID, "expected input's user_id to be passed to loginDB")

	assert.EqualValues(t, userDBMock.ResponseName, resp.Name, "expected returned name to match response from userDB")
	assert.EqualValues(t, AuthDBMock.ResponseAccessToken, resp.Token, "expected token to match response from loginDB")
}

package auth_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	modelauth "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/googleauth"
)

type MockGoogleAuthDB struct {
	inputToken          string
	ResponseUserDetails googleauth.UserDetails
	ResponseError       error
}

func (mock *MockGoogleAuthDB) GetUserDetails(token string) (googleauth.UserDetails, error) {
	mock.inputToken = token

	return mock.ResponseUserDetails, mock.ResponseError
}

func TestGoogleLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError    error
		Input            *modelauth.GoogleLoginRequest
		GoogleAuthDBMock *MockGoogleAuthDB
		UserDBMock       *MockUserDB
		RegisterDBMock   *MockRegisterDB
		AuthDBMock       *MockAuthDB
	}{
		{
			ExpectedError: modelauth.ErrGoogleLoginTokenNotSet,
			Input: &modelauth.GoogleLoginRequest{
				Token: "",
			},
			GoogleAuthDBMock: &MockGoogleAuthDB{},
			UserDBMock:       &MockUserDB{},
			RegisterDBMock:   &MockRegisterDB{},
			AuthDBMock:       &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrGoogleLoginAuth,
			Input: &modelauth.GoogleLoginRequest{
				Token: "token",
			},
			GoogleAuthDBMock: &MockGoogleAuthDB{
				ResponseError: errors.New("test error"),
			},
			UserDBMock:     &MockUserDB{},
			RegisterDBMock: &MockRegisterDB{},
			AuthDBMock:     &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrGoogleLoginGetUser,
			Input: &modelauth.GoogleLoginRequest{
				Token: "token",
			},
			GoogleAuthDBMock: &MockGoogleAuthDB{
				ResponseUserDetails: googleauth.UserDetails{
					Email: "email",
					Name:  "name",
				},
			},
			UserDBMock: &MockUserDB{
				ResponseError: errors.New("test error"),
			},
			RegisterDBMock: &MockRegisterDB{},
			AuthDBMock:     &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrGoogleLoginCreateUser,
			Input: &modelauth.GoogleLoginRequest{
				Token: "token",
			},
			GoogleAuthDBMock: &MockGoogleAuthDB{
				ResponseUserDetails: googleauth.UserDetails{
					Email: "email",
					Name:  "name",
				},
			},
			UserDBMock: &MockUserDB{
				ResponseError: db.ErrGetUserNoRows,
			},
			RegisterDBMock: &MockRegisterDB{
				ResponseError: errors.New("test"),
			},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: modelauth.ErrLoginCreateToken,
			Input: &modelauth.GoogleLoginRequest{
				Token: "token",
			},
			GoogleAuthDBMock: &MockGoogleAuthDB{
				ResponseUserDetails: googleauth.UserDetails{
					Email: "email",
					Name:  "name",
				},
			},
			UserDBMock: &MockUserDB{
				ResponseError: db.ErrGetUserNoRows,
			},
			RegisterDBMock: &MockRegisterDB{},
			AuthDBMock: &MockAuthDB{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		loginModel := modelauth.GoogleLoginModel{
			GoogleAuthDB: test.GoogleAuthDBMock,
			UserDB:       test.UserDBMock,
			RegisterDB:   test.RegisterDBMock,
			AuthDB:       test.AuthDBMock,
		}

		resp, err := loginModel.Do(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		if resp != nil {
			t.Fatalf("expected resp to be nil got %v instead", resp)
		}
	}
}

func TestGoogleLoginSuccess(t *testing.T) {
	t.Parallel()

	googleAuthDBMock := &MockGoogleAuthDB{
		ResponseUserDetails: googleauth.UserDetails{
			Email: "email",
			Name:  "name",
		},
	}

	userDBMock := &MockUserDB{
		ResponseError: nil,
		ResponseName:  googleAuthDBMock.ResponseUserDetails.Name,
	}

	registerDBMock := &MockRegisterDB{
		ResponseError: nil,
		UserID:        googleAuthDBMock.ResponseUserDetails.Email,
		Name:          googleAuthDBMock.ResponseUserDetails.Name,
	}

	authDBMock := &MockAuthDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	googleLoginModel := modelauth.GoogleLoginModel{
		GoogleAuthDB: googleAuthDBMock,
		UserDB:       userDBMock,
		RegisterDB:   registerDBMock,
		AuthDB:       authDBMock,
	}

	input := &modelauth.GoogleLoginRequest{
		Token: "token",
	}

	resp, err := googleLoginModel.Do(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.Token, googleAuthDBMock.inputToken, "expected input token to be passed to GoogleAuthDB")

	assert.EqualValues(t,
		googleAuthDBMock.ResponseUserDetails.Email,
		userDBMock.UserID,
		"expected returned email from googleAuthDB to be passed to userDB as id")
	assert.EqualValues(t,
		googleAuthDBMock.ResponseUserDetails.Email,
		authDBMock.UserID,
		"expected returned email from googleAuthDB to be passed to authDB as id")

	assert.EqualValues(t,
		googleAuthDBMock.ResponseUserDetails.Name,
		resp.Name,
		"expected returned name from googleAuthDB to be passed to resp")
	assert.EqualValues(t,
		authDBMock.ResponseAccessToken,
		resp.Token,
		"expected returned token from authDB to be passed to resp")
}

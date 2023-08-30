package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type MockRegisterDB struct {
	UserID        string
	Name          string
	ResponseError error
}

func (mock *MockRegisterDB) CreateUser(userID string, Name string) error {
	mock.UserID = userID
	mock.Name = Name
	return mock.ResponseError
}

type MockLoginDB struct {
	UserID              string
	ResponseAccessToken login.AccessToken
	ResponseError       error
}

func (mock *MockLoginDB) CreateToken(userID string) (login.AccessToken, error) {
	mock.UserID = userID
	return mock.ResponseAccessToken, mock.ResponseError
}

func TestRegisterErrorHandling(t *testing.T) {
	testCases := []struct {
		ExpectedError  error
		Input          model.RegistrationData
		RegisterDBMock MockRegisterDB
		LoginDBMock    MockLoginDB
	}{
		{
			ExpectedError: model.ErrRegisterUserIDNotSet,
			Input: model.RegistrationData{
				UserID: "",
				Name:   "name",
			},
			RegisterDBMock: MockRegisterDB{},
			LoginDBMock:    MockLoginDB{},
		},

		{
			ExpectedError: model.ErrRegisterNameNotSet,
			Input: model.RegistrationData{
				UserID: "userID",
				Name:   "",
			},
			RegisterDBMock: MockRegisterDB{},
			LoginDBMock:    MockLoginDB{},
		},
		{
			ExpectedError: model.ErrRegisterDuplicate,
			Input: model.RegistrationData{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: MockRegisterDB{
				ResponseError: db.ErrRegisterInsertCount,
			},
			LoginDBMock: MockLoginDB{},
		},
		{
			ExpectedError: model.ErrRegisterCreateUserFailed,
			Input: model.RegistrationData{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: MockRegisterDB{
				ResponseError: errors.New("db error"),
			},
			LoginDBMock: MockLoginDB{},
		},
		{
			ExpectedError: model.ErrRegisterCreateAccessToken,
			Input: model.RegistrationData{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: MockRegisterDB{
				ResponseError: nil,
			},
			LoginDBMock: MockLoginDB{
				ResponseError: errors.New("db error"),
			},
		},
	}

	for _, test := range testCases {
		client := model.Client{
			RegisterDB: &test.RegisterDBMock,
			LoginDB:    &test.LoginDBMock,
		}

		token, err := client.Register(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		assert.EqualValues(t, "", token, "expected token to be empty")
	}
}

func TestRegisterSuccess(t *testing.T) {
	registerDBMock := &MockRegisterDB{
		ResponseError: nil,
	}

	loginDBMock := &MockLoginDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	client := model.Client{
		RegisterDB: registerDBMock,
		LoginDB:    loginDBMock,
	}

	input := model.RegistrationData{
		UserID: "userID",
		Name:   "name",
	}

	token, err := client.Register(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.UserID, registerDBMock.UserID, "expected input user_id to be passed to registerDB")
	assert.EqualValues(t, input.Name, registerDBMock.Name, "expected input name to be passed to registerDB")

	assert.EqualValues(t, loginDBMock.UserID, input.UserID, "expected input's user_id to be passed to loginDB")

	assert.EqualValues(t, loginDBMock.ResponseAccessToken, token, "expected token to match response from loginDB")
}

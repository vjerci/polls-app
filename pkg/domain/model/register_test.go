package model_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type MockRegisterDB struct {
	UserID        string
	Name          string
	ResponseError error
}

func (mock *MockRegisterDB) CreateUser(userID string, name string) error {
	mock.UserID = userID
	mock.Name = name

	return mock.ResponseError
}

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError  error
		Input          *model.RegisterRequest
		RegisterDBMock *MockRegisterDB
		AuthDBMock     *MockAuthDB
	}{
		{
			ExpectedError: model.ErrRegisterUserIDNotSet,
			Input: &model.RegisterRequest{
				UserID: "",
				Name:   "name",
			},
			RegisterDBMock: &MockRegisterDB{},
			AuthDBMock:     &MockAuthDB{},
		},

		{
			ExpectedError: model.ErrRegisterNameNotSet,
			Input: &model.RegisterRequest{
				UserID: "userID",
				Name:   "",
			},
			RegisterDBMock: &MockRegisterDB{},
			AuthDBMock:     &MockAuthDB{},
		},
		{
			ExpectedError: model.ErrRegisterDuplicate,
			Input: &model.RegisterRequest{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: &MockRegisterDB{
				ResponseError: db.ErrCreateUserInsertCount,
			},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: model.ErrRegisterDB,
			Input: &model.RegisterRequest{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: &MockRegisterDB{
				ResponseError: errors.New("db error"),
			},
			AuthDBMock: &MockAuthDB{},
		},
		{
			ExpectedError: model.ErrRegisterCreateAccessToken,
			Input: &model.RegisterRequest{
				UserID: "userID",
				Name:   "name",
			},
			RegisterDBMock: &MockRegisterDB{
				ResponseError: nil,
			},
			AuthDBMock: &MockAuthDB{
				ResponseError: errors.New("db error"),
			},
		},
	}

	for _, test := range testCases {
		registerModel := model.RegisterModel{
			RegisterDB: test.RegisterDBMock,
			AuthDB:     test.AuthDBMock,
		}

		token, err := registerModel.Do(test.Input)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}

		assert.EqualValues(t, "", token, "expected token to be empty")
	}
}

func TestRegisterSuccess(t *testing.T) {
	t.Parallel()

	registerDBMock := &MockRegisterDB{
		ResponseError: nil,
	}

	authDBMock := &MockAuthDB{
		ResponseError:       nil,
		ResponseAccessToken: "testToken",
	}

	registerModel := model.RegisterModel{
		RegisterDB: registerDBMock,
		AuthDB:     authDBMock,
	}

	input := &model.RegisterRequest{
		UserID: "userID",
		Name:   "name",
	}

	token, err := registerModel.Do(input)

	if err != nil {
		t.Fatalf(`expected no err but got "%s" instead`, err)
	}

	assert.EqualValues(t, input.UserID, registerDBMock.UserID, "expected input user_id to be passed to registerDB")
	assert.EqualValues(t, input.Name, registerDBMock.Name, "expected input name to be passed to registerDB")

	assert.EqualValues(t, authDBMock.UserID, input.UserID, "expected input's user_id to be passed to loginDB")

	assert.EqualValues(t, authDBMock.ResponseAccessToken, token, "expected token to match response from loginDB")
}

package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
)

type RegisterRepository interface {
	CreateUser(userID string, Name string) error
}

type RegisterRequest struct {
	UserID string
	Name   string
}

var ErrRegisterUserIDNotSet = errors.New("user_id is not set")
var ErrRegisterNameNotSet = errors.New("name is not set")

var ErrRegisterDuplicate = errors.New("tried to register user that is already registered")

var ErrRegisterDB = errors.New("failed to create user")
var ErrRegisterCreateAccessToken = errors.New("failed to create user")

func (client *Client) Register(data *RegisterRequest) (auth.AccessToken, error) {
	if data.UserID == "" {
		return "", ErrRegisterUserIDNotSet
	}

	if data.Name == "" {
		return "", ErrRegisterNameNotSet
	}

	err := client.RegisterDB.CreateUser(data.UserID, data.Name)
	if err != nil {
		if errors.Is(err, db.ErrCreateUserInsertCount) {
			return "", errors.Join(ErrRegisterDuplicate, err)
		}

		return "", errors.Join(ErrRegisterDB, err)
	}

	token, err := client.AuthDB.CreateToken(data.UserID)
	if err != nil {
		return "", errors.Join(ErrRegisterCreateAccessToken, err)
	}

	return token, nil
}

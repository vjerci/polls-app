package model

import (
	"errors"
	"fmt"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type RegisterRepository interface {
	CreateUser(userID string, groupID string, Name string) error
}

type RegistrationData struct {
	UserID  string `json:"user_id"`
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
}

var ErrRegisterUserIDNotSet = errors.New("user_id is not set")
var ErrRegisterGroupIDNotSet = errors.New("group_id is not set")
var ErrRegisterNameNotSet = errors.New("name is not set")

var ErrRegisterDuplicate = errors.New("tried to register user that is already registered")

var ErrRegisterCreateUserFailed = errors.New("failed to create user")
var ErrRegisterCreateAccessToken = errors.New("failed to create user")

func (client *Client) Register(data RegistrationData) (login.AccessToken, error) {
	if data.UserID == "" {
		return "", ErrRegisterUserIDNotSet
	}

	if data.GroupID == "" {
		return "", ErrRegisterGroupIDNotSet
	}

	if data.Name == "" {
		return "", ErrRegisterNameNotSet
	}

	err := client.RegisterDB.CreateUser(data.UserID, data.GroupID, data.Name)
	if err != nil {
		if errors.Is(err, db.ErrRegisterInsertCount) {
			return "", fmt.Errorf("%w: %w", ErrRegisterDuplicate, err)
		}

		return "", fmt.Errorf("%w: %w", ErrRegisterCreateUserFailed, err)
	}

	token, err := client.LoginDB.CreateToken(data.UserID, data.GroupID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRegisterCreateAccessToken, err)
	}

	return token, nil
}

package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
)

type LoginRequest struct {
	UserID string
}

type LoginResponse struct {
	Token string
	Name  string
}

type LoginRepository interface {
	CreateToken(userID string) (login.AccessToken, error)
}

var ErrLoginUserIDNotSet = errors.New("user id not set")
var ErrLoginUserNotFound = errors.New("user with given id does not exist")
var ErrLoginUserDB = errors.New("getting user failed")
var ErrLoginCreateToken = errors.New("create token failed")

func (client *Client) Login(data *LoginRequest) (*LoginResponse, error) {
	if data.UserID == "" {
		return nil, ErrLoginUserIDNotSet
	}

	name, err := client.UserDB.GetUser(data.UserID)
	if err != nil {
		if errors.Is(err, db.ErrGetUserNoRows) {
			return nil, ErrLoginUserNotFound
		}

		return nil, errors.Join(ErrLoginUserDB, err)
	}

	token, err := client.LoginDB.CreateToken(data.UserID)
	if err != nil {
		return nil, errors.Join(ErrLoginCreateToken, err)
	}

	return &LoginResponse{
		Token: string(token),
		Name:  name,
	}, nil
}

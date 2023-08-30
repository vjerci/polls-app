package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type LoginRepository interface {
	CreateToken(userID string) (login.AccessToken, error)
}

type LoginData struct {
	UserID string `json:"user_id"`
}

var ErrLoginUserIDNotSet = errors.New("user id not set")
var ErrLoginUserNotFound = errors.New("user with given id does not exist")

func (client *Client) Login(data LoginData) (login.AccessToken, error) {
	// if data.UserID == "" {
	// 	return "", ErrLoginUserIDNotSet
	// }

	// name, err := client.UserDB.GetUser(data.UserID)
	// if err != nil {
	// 	if errors.Is(err, db.ErrUserNoRows) {
	// 		return "", ErrLoginUserNotFound
	// 	}
	// }

	token, err := client.LoginDB.CreateToken(data.UserID)
	if err != nil {
		return "", errors.Join(ErrLoginUserIDNotSet, err)
	}

	return token, nil
}

package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
)

type UserRepository interface {
	GetUser(userID string) (name string, err error)
}

type UserModel struct {
	UserDB UserRepository
}

var ErrGetUserUserIDEmpty = errors.New("user must not be empty")
var ErrGetUserUserNotFound = errors.New("user not found")

var ErrGetUserUserDB = errors.New("error while getting user")

func (model *UserModel) GetUser(userID string) (string, error) {
	if userID == "" {
		return "", ErrGetUserUserIDEmpty
	}

	name, err := model.UserDB.GetUser(userID)
	if err != nil {
		if errors.Is(err, db.ErrGetUserNoRows) {
			return "", ErrGetUserUserNotFound
		}

		return "", errors.Join(ErrGetUserUserDB, err)
	}

	return name, nil
}

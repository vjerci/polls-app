package auth

import (
	"errors"

	"github.com/vjerci/polls-app/pkg/domain/db"
)

type LoginRequest struct {
	UserID string
}

type LoginResponse struct {
	Token string
	Name  string
}

type LoginModel struct {
	AuthDB AuthRepository
	UserDB UserRepository
}

var ErrLoginUserIDNotSet = errors.New("user id not set")
var ErrLoginUserNotFound = errors.New("user with given id does not exist")
var ErrLoginUserDB = errors.New("getting user failed")
var ErrLoginCreateToken = errors.New("create token failed")

func (model *LoginModel) Do(data *LoginRequest) (resp *LoginResponse, err error) {
	if data.UserID == "" {
		return nil, ErrLoginUserIDNotSet
	}

	name, err := model.UserDB.GetUser(data.UserID)
	if err != nil {
		if errors.Is(err, db.ErrGetUserNoRows) {
			return nil, ErrLoginUserNotFound
		}

		return nil, errors.Join(ErrLoginUserDB, err)
	}

	token, err := model.AuthDB.CreateToken(data.UserID)
	if err != nil {
		return nil, errors.Join(ErrLoginCreateToken, err)
	}

	return &LoginResponse{
		Token: string(token),
		Name:  name,
	}, nil
}

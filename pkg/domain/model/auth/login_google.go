package auth

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/googleauth"
)

type LoginGoogleRequest struct {
	Token string
}

type LoginGoogleResponse struct {
	Token string
	Name  string
}

type GoogleAuthRepository interface {
	GetUserDetails(token string) (googleauth.UserDetails, error)
}

type LoginGoogleModel struct {
	GoogleAuthDB GoogleAuthRepository
	AuthDB       AuthRepository
	UserDB       UserRepository
	RegisterDB   RegisterRepository
}

var ErrLoginGoogleTokenNotSet = errors.New("google token not set")
var ErrLoginGoogleAuth = errors.New("failed to authenticate with google")
var ErrLoginGoogleGetUser = errors.New("failed to get user")
var ErrLoginGoogleCreateUser = errors.New("failed to create user")

func (model *LoginGoogleModel) Do(data *LoginGoogleRequest) (resp *LoginGoogleResponse, err error) {
	if data.Token == "" {
		return nil, ErrLoginGoogleTokenNotSet
	}

	userDetails, err := model.GoogleAuthDB.GetUserDetails(data.Token)
	if err != nil {
		return nil, errors.Join(ErrLoginGoogleAuth, err)
	}

	// try to get user
	name, err := model.UserDB.GetUser(userDetails.Email)
	if err != nil && !errors.Is(err, db.ErrGetUserNoRows) {
		return nil, errors.Join(ErrLoginGoogleGetUser, err)
	}

	// if user doesn't exist create it
	if errors.Is(err, db.ErrGetUserNoRows) {
		err := model.RegisterDB.CreateUser(userDetails.Email, userDetails.Name)
		if err != nil {
			return nil, errors.Join(ErrLoginGoogleCreateUser, err)
		}
	}

	token, err := model.AuthDB.CreateToken(userDetails.Email)
	if err != nil {
		return nil, errors.Join(ErrLoginCreateToken, err)
	}

	return &LoginGoogleResponse{
		Token: string(token),
		Name:  name,
	}, nil
}

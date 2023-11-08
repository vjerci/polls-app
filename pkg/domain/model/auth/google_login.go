package auth

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/googleauth"
)

type GoogleLoginRequest struct {
	Token string
}

type GoogleLoginResponse struct {
	Token string
	Name  string
}

type GoogleAuthRepository interface {
	GetUserDetails(token string) (googleauth.UserDetails, error)
}

type GoogleLoginModel struct {
	GoogleAuthDB GoogleAuthRepository
	AuthDB       AuthRepository
	UserDB       UserRepository
	RegisterDB   RegisterRepository
}

var ErrGoogleLoginTokenNotSet = errors.New("google token not set")
var ErrGoogleLoginAuth = errors.New("failed to authenticate with google")
var ErrGoogleLoginGetUser = errors.New("failed to get user")
var ErrGoogleLoginCreateUser = errors.New("failed to create user")

func (model *GoogleLoginModel) Do(data *GoogleLoginRequest) (resp *GoogleLoginResponse, err error) {
	if data.Token == "" {
		return nil, ErrGoogleLoginTokenNotSet
	}

	userDetails, err := model.GoogleAuthDB.GetUserDetails(data.Token)
	if err != nil {
		return nil, errors.Join(ErrGoogleLoginAuth, err)
	}

	// try to get user
	name, err := model.UserDB.GetUser(userDetails.Email)
	if err != nil && !errors.Is(err, db.ErrGetUserNoRows) {
		return nil, errors.Join(ErrGoogleLoginGetUser, err)
	}

	// if user doesn't exist create it
	if errors.Is(err, db.ErrGetUserNoRows) {
		err := model.RegisterDB.CreateUser(userDetails.Email, userDetails.Name)
		if err != nil {
			return nil, errors.Join(ErrGoogleLoginCreateUser, err)
		}
	}

	token, err := model.AuthDB.CreateToken(userDetails.Email)
	if err != nil {
		return nil, errors.Join(ErrLoginCreateToken, err)
	}

	return &GoogleLoginResponse{
		Token: string(token),
		Name:  name,
	}, nil
}

package googleauth

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type Client struct {
	ClientID string
}

func New(clientID string) *Client {
	return &Client{
		ClientID: clientID,
	}
}

type UserDetails struct {
	Email string
	Name  string
}

var ErrGoogleLoginValidation = errors.New("failed to validate google token")
var ErrGoogleLoginEmailNotSet = errors.New("didn't get an email response from google")
var ErrGoogleLoginEmailNotString = errors.New("email response from google is not a string")
var ErrGoogleLoginUserNameNotSet = errors.New("didn't get an name response from google")
var ErrGoogleLoginUserNameNotString = errors.New("name response from google is not a string")

func (client *Client) GetUserDetails(token string) (UserDetails, error) {
	googleToken, err := idtoken.Validate(context.Background(), token, client.ClientID)
	if err != nil {
		return UserDetails{}, errors.Join(ErrGoogleLoginValidation, err)
	}

	userEmailInterface, emailExists := googleToken.Claims["email"]
	if !emailExists {
		return UserDetails{}, ErrGoogleLoginEmailNotSet
	}

	userEmail, isString := userEmailInterface.(string)
	if !isString {
		return UserDetails{}, ErrGoogleLoginEmailNotString
	}

	userNameInterface, nameExists := googleToken.Claims["name"]
	if !nameExists {
		return UserDetails{}, ErrGoogleLoginUserNameNotSet
	}

	userName, isString := userNameInterface.(string)
	if !isString {
		return UserDetails{}, ErrGoogleLoginUserNameNotString
	}

	return UserDetails{
		Email: userEmail,
		Name:  userName,
	}, nil
}

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

var ErrLoginGoogleValidation = errors.New("failed to validate google token")
var ErrLoginGoogleEmailNotSet = errors.New("didn't get an email response from google")
var ErrLoginGooglEmailNotString = errors.New("email response from google is not a string")
var ErrLoginGoogleUserNameNotSet = errors.New("didn't get an name response from google")
var ErrLoginGooglUserNameNotString = errors.New("name response from google is not a string")

func (client *Client) GetUserDetails(token string) (UserDetails, error) {
	googleToken, err := idtoken.Validate(context.Background(), token, client.ClientID)
	if err != nil {
		return UserDetails{}, errors.Join(ErrLoginGoogleValidation, err)
	}

	userEmailInterface, emailExists := googleToken.Claims["email"]
	if !emailExists {
		return UserDetails{}, ErrLoginGoogleEmailNotSet
	}

	userEmail, isString := userEmailInterface.(string)
	if !isString {
		return UserDetails{}, ErrLoginGooglEmailNotString
	}

	userNameInterface, nameExists := googleToken.Claims["name"]
	if !nameExists {
		return UserDetails{}, ErrLoginGoogleEmailNotSet
	}

	userName, isString := userNameInterface.(string)
	if !isString {
		return UserDetails{}, ErrLoginGooglUserNameNotString
	}

	return UserDetails{
		Email: userEmail,
		Name:  userName,
	}, nil
}

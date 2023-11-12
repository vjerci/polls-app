package middleware

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type Client struct {
	AuthRepo AuthRepo
	UserRepo UserRepo
}

type AuthRepo interface {
	Decode(input auth.AccessToken) (userID string, err error)
}

type UserRepo interface {
	GetUser(userID string) (name string, err error)
}

func (client *Client) WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		header := echoContext.Request().Header.Get("Authorization")
		if header == "" {
			return schema.ErrAuthMissing
		}

		token, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			return schema.ErrAuthInvalid
		}

		userID, err := client.AuthRepo.Decode(auth.AccessToken(token))
		if err != nil {
			return schema.ErrAuthInvalid
		}

		_, err = client.UserRepo.GetUser(userID)
		if err != nil {
			if errors.Is(err, authmodel.ErrGetUserUserNotFound) {
				return schema.ErrAuthUserDoesNotExist
			}

			return err
		}

		echoContext.Set("userID", userID)

		return next(echoContext)
	}
}

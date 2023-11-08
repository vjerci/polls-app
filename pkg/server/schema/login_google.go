package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
)

type LoginGoogleRequest struct {
	Token string `json:"token"`
}

func (req *LoginGoogleRequest) ToModel() *auth.LoginGoogleRequest {
	return &auth.LoginGoogleRequest{
		Token: req.Token,
	}
}

type LoginGoogleResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

type LoginGoogleSchemaMap struct{}

func (mapper *LoginGoogleSchemaMap) Response(input *auth.LoginGoogleResponse) *LoginGoogleResponse {
	return &LoginGoogleResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

var ErrLoginGoogleTokenNotSet = &echo.HTTPError{
	Message:  `empty token submitted`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrLoginGoogleJSONDecode = &echo.HTTPError{
	Message:  "failed to decode login with google json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrLoginGoogleModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *LoginGoogleSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, auth.ErrLoginGoogleTokenNotSet) {
		return ErrLoginGoogleTokenNotSet.WithInternal(err)
	}

	return ErrLoginGoogleModel.WithInternal(err)
}

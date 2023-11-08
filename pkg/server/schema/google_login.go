package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
)

type GoogleLoginRequest struct {
	Token string `json:"token"`
}

func (req *GoogleLoginRequest) ToModel() *auth.GoogleLoginRequest {
	return &auth.GoogleLoginRequest{
		Token: req.Token,
	}
}

type GoogleLoginResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

type GoogleLoginSchemaMap struct{}

func (mapper *GoogleLoginSchemaMap) Response(input *auth.GoogleLoginResponse) *GoogleLoginResponse {
	return &GoogleLoginResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

var ErrGoogleLoginTokenNotSet = &echo.HTTPError{
	Message:  `empty token submitted`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrGoogleLoginJSONDecode = &echo.HTTPError{
	Message:  "failed to decode login with google json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrGoogleLoginModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *GoogleLoginSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, auth.ErrGoogleLoginTokenNotSet) {
		return ErrGoogleLoginTokenNotSet.WithInternal(err)
	}

	return ErrGoogleLoginModel.WithInternal(err)
}

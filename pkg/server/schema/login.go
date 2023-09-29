package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type LoginRequest struct {
	UserID string `json:"user_id"`
}

func (req *LoginRequest) ToModel() *model.LoginRequest {
	return &model.LoginRequest{
		UserID: req.UserID,
	}
}

type LoginSchemaMap struct{}

type LoginResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func (mapper *LoginSchemaMap) Response(input *model.LoginResponse) *LoginResponse {
	return &LoginResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

var ErrLoginUserDoesNotExist = &echo.HTTPError{
	Message:  `user with given "user_id" does not exist`,
	Code:     http.StatusNotFound,
	Internal: nil,
}
var ErrLoginUserIDNotSet = &echo.HTTPError{
	Message:  `input field "user_id" not set`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrLoginJSONDecode = &echo.HTTPError{
	Message:  "failed to decode login json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrLoginModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *LoginSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrLoginUserIDNotSet) {
		return ErrLoginUserIDNotSet.WithInternal(err)
	}

	if errors.Is(err, model.ErrLoginUserNotFound) {
		return ErrLoginUserDoesNotExist.WithInternal(err)
	}

	return ErrLoginModel.WithInternal(err)
}

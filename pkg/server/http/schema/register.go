package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
)

type RegisterRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (req *RegisterRequest) ToModel() *authmodel.RegisterRequest {
	return &authmodel.RegisterRequest{
		UserID: req.UserID,
		Name:   req.Name,
	}
}

type RegisterSchemaMap struct{}

func (mapper *RegisterSchemaMap) Response(input auth.AccessToken) string {
	return string(input)
}

var ErrRegisterUserIDNotSet = &echo.HTTPError{
	Message:  `field "user_id" is not set`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrRegisterNameNotSet = &echo.HTTPError{
	Message:  `field "name" is not set`,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrRegisterUserDuplicate = &echo.HTTPError{
	Message:  `user with given "user_id" already registered`,
	Code:     http.StatusConflict,
	Internal: nil,
}

var ErrRegisterModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

var ErrRegisterJSONDecode = &echo.HTTPError{
	Message:  "failed to decode registration json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

func (mapper *RegisterSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, authmodel.ErrRegisterUserIDNotSet) {
		return ErrRegisterUserIDNotSet.WithInternal(err)
	}

	if errors.Is(err, authmodel.ErrRegisterNameNotSet) {
		return ErrRegisterNameNotSet.WithInternal(err)
	}

	if errors.Is(err, authmodel.ErrRegisterDuplicate) {
		return ErrRegisterUserDuplicate.WithInternal(err)
	}

	return ErrRegisterModel.WithInternal(err)
}

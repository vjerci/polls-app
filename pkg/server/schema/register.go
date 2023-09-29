package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
)

type RegisterRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (req *RegisterRequest) ToModel() *model.RegisterRequest {
	return &model.RegisterRequest{
		UserID: req.UserID,
		Name:   req.Name,
	}
}

type RegisterSchemaMap struct{}

func (mapper *RegisterSchemaMap) Response(input auth.AccessToken) string {
	return string(input)
}

var ErrRegisterUserIDNotSet = &echo.HTTPError{
	Message:  model.ErrRegisterUserIDNotSet,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrRegisterNameNotSet = &echo.HTTPError{
	Message:  model.ErrRegisterNameNotSet,
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrRegisterUserDuplicate = &echo.HTTPError{
	Message:  model.ErrRegisterDuplicate,
	Code:     http.StatusConflict,
	Internal: nil,
}

var ErrRegisterModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

var ErrRegisterJSONDecode = &echo.HTTPError{
	Message:  errors.New("failed to decode registration json body"),
	Code:     http.StatusBadRequest,
	Internal: nil,
}

func (mapper *RegisterSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrRegisterUserIDNotSet) {
		return ErrRegisterUserIDNotSet.WithInternal(err)
	}

	if errors.Is(err, model.ErrRegisterNameNotSet) {
		return ErrRegisterNameNotSet.WithInternal(err)
	}

	if errors.Is(err, model.ErrRegisterDuplicate) {
		return ErrRegisterUserDuplicate.WithInternal(err)
	}

	return ErrRegisterModel.WithInternal(err)
}

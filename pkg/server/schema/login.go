package schema

import (
	"errors"
	"net/http"

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

type LoginResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func (mapper *Map) LoginResponse(input *model.LoginResponse) *LoginResponse {
	return &LoginResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

var ErrLoginUserDoesNotExist = &UserVisibleError{
	Err:    model.ErrLoginUserNotFound,
	Status: http.StatusNotFound,
}
var ErrLoginUserIDNotSet = &UserVisibleError{
	Err:    model.ErrLoginUserIDNotSet,
	Status: http.StatusBadRequest,
}
var handledLoginErrors = []*UserVisibleError{
	ErrLoginUserDoesNotExist,
	ErrLoginUserIDNotSet,
}

var ErrLoginJSONDecode = &UserVisibleError{
	Err:    errors.New("failed to decode login json body"),
	Status: http.StatusBadRequest,
}

var ErrLoginModel = errors.New("model failed to login")

func (mapper *Map) LoginError(err error) error {
	for _, targetError := range handledLoginErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrLoginModel, err)
}

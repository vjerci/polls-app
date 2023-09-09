package schema

import (
	"errors"
	"net/http"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
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

func (mapper *Map) RegisterResponse(input login.AccessToken) string {
	return string(input)
}

var ErrRegisterUserIDNotSet = &UserVisibleError{
	Err:    model.ErrRegisterUserIDNotSet,
	Status: http.StatusBadRequest,
}
var ErrRegisterNameNotSet = &UserVisibleError{
	Err:    model.ErrRegisterNameNotSet,
	Status: http.StatusBadRequest,
}
var ErrRegisterUserDuplicate = &UserVisibleError{
	Err:    model.ErrRegisterDuplicate,
	Status: http.StatusConflict,
}
var handledRegisterErrors = []*UserVisibleError{
	ErrRegisterUserIDNotSet,
	ErrRegisterNameNotSet,
	ErrRegisterUserDuplicate,
}

var ErrRegisterModel = errors.New("model failed to register")

var ErrRegisterJSONDecode = &UserVisibleError{
	Err:    errors.New("failed to decode registration json body"),
	Status: http.StatusBadRequest,
}

func (mapper *Map) RegisterError(err error) error {
	for _, targetError := range handledRegisterErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrRegisterModel, err)
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type RegisterRepository interface {
	Register(model.RegistrationData) (login.AccessToken, error)
}

var ErrRegisterJsonDecode = &ErrorUserVisible{Err: errors.New("failed to decode registration json body"), Status: http.StatusBadRequest}

var ErrRegisterUserIDNotSet = &ErrorUserVisible{Err: model.ErrRegisterUserIDNotSet, Status: http.StatusBadRequest}
var ErrRegisterNameNotSet = &ErrorUserVisible{Err: model.ErrRegisterNameNotSet, Status: http.StatusBadRequest}
var ErrRegisterUserDuplicate = &ErrorUserVisible{Err: model.ErrRegisterDuplicate, Status: http.StatusConflict}

var ErrRegisterModel = errors.New("model failed to register")

func (api *FactoryImplementation) Register(repo RegisterRepository) echo.HandlerFunc {
	return func(c echo.Context) error {

		var data model.RegistrationData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return &ErrorUserVisible{Err: errors.Join(ErrRegisterJsonDecode, err), Status: ErrRegisterJsonDecode.Status}
		}

		var handledRepositoryErrors = []*ErrorUserVisible{ErrRegisterUserIDNotSet, ErrRegisterNameNotSet, ErrRegisterUserDuplicate}

		token, err := repo.Register(data)
		if err != nil {
			for _, targetError := range handledRepositoryErrors {
				if errors.Is(err, targetError.Err) {
					return targetError
				}
			}

			return errors.Join(ErrRegisterModel, err)
		}

		return c.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Data:    token,
		})
	}
}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type RegisterRepository interface {
	Register(model.RegistrationData) (login.AccessToken, error)
}

var ErrRegisterJsonDecode = &ErrorClientVisible{Err: errors.New("failed to decode registration json body"), Status: http.StatusBadRequest}

var ErrRegisterUserIDNotSet = &ErrorClientVisible{Err: model.ErrRegisterUserIDNotSet, Status: http.StatusBadRequest}
var ErrRegisterGroupIDNotSet = &ErrorClientVisible{Err: model.ErrRegisterGroupIDNotSet, Status: http.StatusBadRequest}
var ErrRegisterNameNotSet = &ErrorClientVisible{Err: model.ErrRegisterNameNotSet, Status: http.StatusBadRequest}
var ErrRegisterUserDuplicate = &ErrorClientVisible{Err: model.ErrRegisterDuplicate, Status: http.StatusConflict}

var ErrRegisterModel = errors.New("model failed to register")

func (api *FactoryImplementation) Register(repo RegisterRepository) echo.HandlerFunc {
	var handledRepositoryErrors = []*ErrorClientVisible{ErrRegisterUserIDNotSet, ErrRegisterGroupIDNotSet, ErrRegisterNameNotSet, ErrRegisterUserDuplicate}

	return func(c echo.Context) error {
		var data model.RegistrationData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return &ErrorClientVisible{Err: fmt.Errorf("%w: %w", ErrRegisterJsonDecode, err), Status: ErrRegisterJsonDecode.Status}
		}

		token, err := repo.Register(data)
		if err != nil {
			for _, targetError := range handledRepositoryErrors {
				if errors.Is(targetError, err) {
					return targetError
				}
			}

			return fmt.Errorf("%w: %w", ErrRegisterModel, err)
		}

		return c.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Data:    token,
		})
	}
}

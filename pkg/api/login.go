package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

type LoginModel interface {
	Login(data model.LoginData) (accessToken login.AccessToken, name string, err error)
}

var ErrLoginJsonDecode = &ErrorUserVisible{Err: errors.New("failed to decode login json body"), Status: http.StatusBadRequest}
var ErrLoginUserDoesNotExist = &ErrorUserVisible{Err: model.ErrLoginUserNotFound, Status: http.StatusNotFound}
var ErrLoginUserIDNotSet = &ErrorUserVisible{Err: model.ErrLoginUserIDNotSet, Status: http.StatusBadRequest}

var ErrorLoginModel = errors.New("couldn't get user and token from model")

func (factory *FactoryImplementation) Login(repo LoginModel) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data model.LoginData
		err := json.NewDecoder(c.Request().Body).Decode(&data)
		if err != nil {
			return &ErrorUserVisible{Err: errors.Join(ErrLoginJsonDecode, err), Status: ErrLoginJsonDecode.Status}
		}

		var handledRepositoryErrors = []*ErrorUserVisible{ErrLoginUserDoesNotExist, ErrLoginUserIDNotSet}

		token, name, err := repo.Login(data)
		if err != nil {
			for _, targetError := range handledRepositoryErrors {
				if errors.Is(err, targetError.Err) {
					return targetError
				}
			}

			return errors.Join(ErrorLoginModel, err)
		}

		return c.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Data: struct {
				Name        string `json:"name"`
				AccessToken string `json:"access_token"`
			}{
				Name:        name,
				AccessToken: string(token),
			},
		})
	}
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type RegisterModel interface {
	Register(input *model.RegisterRequest) (accessToken login.AccessToken, err error)
}

type RegisterSchemaMap interface {
	RegisterError(err error) error
	RegisterResponse(input login.AccessToken) string
}

func (factory *FactoryImplementation) Register(
	registerModel RegisterModel,
	schemaMap RegisterSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		var data schema.RegisterRequest

		err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
		if err != nil {
			return &schema.UserVisibleError{
				Err:    errors.Join(schema.ErrRegisterJSONDecode, err),
				Status: schema.ErrRegisterJSONDecode.Status,
			}
		}

		token, err := registerModel.Register(data.ToModel())
		if err != nil {
			return schemaMap.RegisterError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.RegisterResponse(token),
			Error:   nil,
		})
	}
}

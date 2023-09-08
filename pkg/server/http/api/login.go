package api

import (
	"encoding/json"
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type LoginModel interface {
	Login(data *model.LoginRequest) (resp *model.LoginResponse, err error)
}

type LoginSchemaMap interface {
	LoginError(err error) error
	LoginResponse(input *model.LoginResponse) *schema.LoginResponse
}

func (factory *FactoryImplementation) Login(
	loginModel LoginModel,
	schemaMap LoginSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		var data schema.LoginRequest

		err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
		if err != nil {
			return &schema.UserVisibleError{
				Err:    errors.Join(schema.ErrLoginJSONDecode, err),
				Status: schema.ErrLoginJSONDecode.Status,
			}
		}

		resp, err := loginModel.Login(data.ToModel())
		if err != nil {
			return schemaMap.LoginError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.LoginResponse(resp),
			Error:   nil,
		})
	}
}

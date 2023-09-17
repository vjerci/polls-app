package api

import (
	"encoding/json"
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type RegisterModel interface {
	Do(input *model.RegisterRequest) (accessToken auth.AccessToken, err error)
}

type RegisterSchemaMap interface {
	ErrorHandler(err error) error
	Response(input auth.AccessToken) string
}

func (client *API) Register(echoContext echo.Context) error {
	var data schema.RegisterRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return &schema.UserVisibleError{
			Err:    errors.Join(schema.ErrRegisterJSONDecode, err),
			Status: schema.ErrRegisterJSONDecode.Status,
		}
	}

	token, err := client.models.Register.Do(data.ToModel())
	if err != nil {
		return client.schemas.Register.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.Register.Response(token),
		Error:   nil,
	})
}

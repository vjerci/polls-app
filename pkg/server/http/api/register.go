package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	modelauth "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type RegisterModel interface {
	Do(input *modelauth.RegisterRequest) (accessToken auth.AccessToken, err error)
}

type RegisterSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input auth.AccessToken) string
}

func (client *API) Register(echoContext echo.Context) error {
	var data schema.RegisterRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrRegisterJSONDecode.WithInternal(err)
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

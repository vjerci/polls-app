package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type LoginGoogleModel interface {
	Do(data *auth.LoginGoogleRequest) (resp *auth.LoginGoogleResponse, err error)
}

type LoginGoogleSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *auth.LoginGoogleResponse) *schema.LoginGoogleResponse
}

func (client *API) LoginGoogle(echoContext echo.Context) error {
	var data schema.LoginGoogleRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrLoginGoogleJSONDecode.WithInternal(err)
	}

	resp, err := client.models.LoginGoogle.Do(data.ToModel())
	if err != nil {
		return client.schemas.Login.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.LoginGoogle.Response(resp),
		Error:   nil,
	})
}

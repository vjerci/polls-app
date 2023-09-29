package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type LoginModel interface {
	Do(data *model.LoginRequest) (resp *model.LoginResponse, err error)
}

type LoginSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *model.LoginResponse) *schema.LoginResponse
}

func (client *API) Login(echoContext echo.Context) error {
	var data schema.LoginRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrLoginJSONDecode.WithInternal(err)
	}

	resp, err := client.models.Login.Do(data.ToModel())
	if err != nil {
		return client.schemas.Login.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.Login.Response(resp),
		Error:   nil,
	})
}

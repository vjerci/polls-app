package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type GoogleLoginModel interface {
	Do(data *auth.GoogleLoginRequest) (resp *auth.GoogleLoginResponse, err error)
}

type GoogleLoginSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *auth.GoogleLoginResponse) *schema.GoogleLoginResponse
}

func (client *API) GoogleLogin(echoContext echo.Context) error {
	var data schema.GoogleLoginRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrGoogleLoginJSONDecode.WithInternal(err)
	}

	resp, err := client.models.GoogleLogin.Do(data.ToModel())
	if err != nil {
		return client.schemas.GoogleLogin.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.GoogleLogin.Response(resp),
		Error:   nil,
	})
}

package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
)

func (client *Client) ErrorHandler(err error, echoContext echo.Context) {
	//nolint:errorlint
	if httpError, ok := err.(*echo.HTTPError); ok {
		message := "internal server error"

		if messageString, ok := httpError.Message.(string); ok {
			message = messageString
		}

		echoContext.Logger().Errorf(`error serving http error "%s":"%s"`, httpError.Message, httpError.Internal)

		err := echoContext.JSON(httpError.Code, api.Response{
			Success: false,
			Error:   &message,
			Data:    nil,
		})

		if err != nil {
			echoContext.Logger().Errorf("failed to serve error response %s", err)
		}
	}

	defaultMessage := "internal server error"

	err = echoContext.JSON(http.StatusInternalServerError, api.Response{
		Success: false,
		Error:   &defaultMessage,
		Data:    nil,
	})

	if err != nil {
		echoContext.Logger().Errorf("failed to serve error response %s", err)
	}
}

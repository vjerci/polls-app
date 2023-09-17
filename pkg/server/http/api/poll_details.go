package api

import (
	"errors"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollDetailsModel interface {
	Get(data *model.PollDetailsRequest) (*model.PollDetailsResponse, error)
}

type PollDetailsSchemaMap interface {
	ErrorHandler(err error) error
	Response(input *model.PollDetailsResponse) *schema.PollDetailsResponse
}

func (client *API) PollDetails(echoContext echo.Context) error {
	userID := echoContext.Get("userID")

	userIDS, ok := userID.(string)
	if !ok {
		return errors.Join(ErrUserIDIsNotString, fmt.Errorf("got %#v for userID", userID))
	}

	resp, err := client.models.PollDetails.Get(&model.PollDetailsRequest{
		PollID: echoContext.Param("id"),
		UserID: userIDS,
	})
	if err != nil {
		return client.schemas.PollDetails.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.PollDetails.Response(resp),
		Error:   nil,
	})
}

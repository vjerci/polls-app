package api

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type PollDetailsModel interface {
	Get(data *poll.DetailsRequest) (*poll.DetailsResponse, error)
}

type PollDetailsSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *poll.DetailsResponse) *schema.PollDetailsResponse
}

func (client *API) PollDetails(echoContext echo.Context) error {
	userID := echoContext.Get("userID")

	userIDS, ok := userID.(string)
	if !ok {
		return ErrUserIDIsNotString.WithInternal(fmt.Errorf("got %#v for userID", userID))
	}

	resp, err := client.models.PollDetails.Get(&poll.DetailsRequest{
		ID:     echoContext.Param("id"),
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

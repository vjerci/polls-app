package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

type PollCreateModel interface {
	Create(data *poll.CreateRequest) (*poll.CreateResponse, error)
}

type PollCreateSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *poll.CreateResponse) *schema.PollCreateResponse
}

func (client *API) PollCreate(echoContext echo.Context) error {
	var data schema.PollCreateRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrPollCreateJSONDecode.WithInternal(err)
	}

	resp, err := client.models.PollCreate.Create(data.ToModel())
	if err != nil {
		return client.schemas.PollCreate.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.PollCreate.Response(resp),
		Error:   nil,
	})
}

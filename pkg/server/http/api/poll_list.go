package api

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollListModel interface {
	Get(data *model.PollListRequest) (*model.PollListResponse, error)
}

type PollListSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *model.PollListResponse) *schema.PollListResponse
}

func (client *API) PollList(echoContext echo.Context) error {
	var data schema.PollListRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrPollListJSONDecode.WithInternal(err)
	}

	resp, err := client.models.PollList.Get(data.ToModel())
	if err != nil {
		return client.schemas.PollList.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.PollList.Response(resp),
		Error:   nil,
	})
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollCreateModel interface {
	Create(data *model.PollCreateRequest) (*model.PollCreateResponse, error)
}

type PollCreateSchemaMap interface {
	ErrorHandler(err error) error
	Response(input *model.PollCreateResponse) *schema.PollCreateResponse
}

func (client *API) PollCreate(echoContext echo.Context) error {
	var data schema.PollCreateRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return &schema.UserVisibleError{
			Err:    errors.Join(schema.ErrPollCreateJSONDecode, err),
			Status: schema.ErrPollCreateJSONDecode.Status,
		}
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

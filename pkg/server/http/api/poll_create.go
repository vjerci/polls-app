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
	CreatePoll(data *model.PollCreateRequest) (*model.PollCreateResponse, error)
}

type PollCreateSchemaMap interface {
	PollCreateError(err error) error
	PollCreateResponse(input *model.PollCreateResponse) *schema.PollCreateResponse
}

func (factory *FactoryImplementation) PollCreate(
	pollCreateModel PollCreateModel,
	schemaMap PollCreateSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		var data schema.PollCreateRequest

		err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
		if err != nil {
			return &schema.UserVisibleError{
				Err:    errors.Join(schema.ErrPollCreateJSONDecode, err),
				Status: schema.ErrPollCreateJSONDecode.Status,
			}
		}

		resp, err := pollCreateModel.CreatePoll(data.ToModel())
		if err != nil {
			return schemaMap.PollCreateError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.PollCreateResponse(resp),
			Error:   nil,
		})
	}
}

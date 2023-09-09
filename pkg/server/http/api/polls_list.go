package api

import (
	"encoding/json"
	"errors"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollListModel interface {
	GetPollList(data *model.PollListRequest) (*model.PollListResponse, error)
}

type PollListSchemaMap interface {
	PollListError(err error) error
	PollListResponse(input *model.PollListResponse) *schema.PollListResponse
}

func (factory *FactoryImplementation) PollList(
	pollListModel PollListModel,
	schemaMap PollListSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		var data schema.PollListRequest

		err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
		if err != nil {
			return &schema.UserVisibleError{
				Err:    errors.Join(schema.ErrPollListJSONDecode, err),
				Status: schema.ErrPollListJSONDecode.Status,
			}
		}

		resp, err := pollListModel.GetPollList(data.ToModel())
		if err != nil {
			return schemaMap.PollListError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.PollListResponse(resp),
			Error:   nil,
		})
	}
}

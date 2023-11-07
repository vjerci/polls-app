package api

import (
	"net/http"
	"strconv"

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
	pageStr := echoContext.QueryParam("page")

	page := 0

	if pageStr != "" {
		var err error

		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return schema.ErrPollListPageNotInt
		}
	}

	resp, err := client.models.PollList.Get(&model.PollListRequest{
		Page: page,
	})
	if err != nil {
		return client.schemas.PollList.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.PollList.Response(resp),
		Error:   nil,
	})
}

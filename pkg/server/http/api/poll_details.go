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
	GetPollDetails(data *model.PollDetailsRequest) (*model.PollDetailsResponse, error)
}

type PollDetailsSchemaMap interface {
	PollDetailsError(err error) error
	PollDetailsResponse(input *model.PollDetailsResponse) *schema.PollDetailsResponse
}

func (factory *FactoryImplementation) PollDetails(
	pollDetailsModel PollDetailsModel,
	schemaMap PollDetailsSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		userID := echoContext.Get("userID")

		userIDS, ok := userID.(string)
		if !ok {
			return errors.Join(ErrUserIDIsNotString, fmt.Errorf("got %#v for userID", userID))
		}

		resp, err := pollDetailsModel.GetPollDetails(&model.PollDetailsRequest{
			PollID: echoContext.Param("id"),
			UserID: userIDS,
		})
		if err != nil {
			return schemaMap.PollDetailsError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.PollDetailsResponse(resp),
			Error:   nil,
		})
	}
}

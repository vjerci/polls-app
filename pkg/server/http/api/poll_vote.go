package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollVoteModel interface {
	PollVote(data *model.PollVoteRequest) (*model.PollVoteResponse, error)
}

type PollVoteSchemaMap interface {
	PollVoteError(err error) error
	PollVoteResponse(*model.PollVoteResponse) *schema.PollVoteResponse
}

func (factory *FactoryImplementation) PollVote(
	pollVoteModel PollVoteModel,
	schemaMap PollVoteSchemaMap,
) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		userID := echoContext.Get("userID")

		userIDS, ok := userID.(string)
		if !ok {
			return errors.Join(ErrUserIDIsNotString, fmt.Errorf("got %#v for userID", userID))
		}

		var data schema.PollVoteRequest

		err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
		if err != nil {
			return &schema.UserVisibleError{
				Err:    errors.Join(schema.ErrPollVoteJSONDecode, err),
				Status: schema.ErrPollVoteJSONDecode.Status,
			}
		}

		resp, err := pollVoteModel.PollVote(&model.PollVoteRequest{
			PollID:   echoContext.Param("id"),
			AnswerID: data.AnswerID,
			UserID:   userIDS,
		})
		if err != nil {
			return schemaMap.PollVoteError(err)
		}

		return echoContext.JSON(http.StatusOK, Response{
			Success: true,
			Data:    schemaMap.PollVoteResponse(resp),
			Error:   nil,
		})
	}
}

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
	Do(data *model.PollVoteRequest) (*model.PollVoteResponse, error)
}

type PollVoteSchemaMap interface {
	ErrorHandler(err error) error
	Response(*model.PollVoteResponse) *schema.PollVoteResponse
}

func (client *API) PollVote(echoContext echo.Context) error {
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

	resp, err := client.models.PollVote.Do(&model.PollVoteRequest{
		PollID:   echoContext.Param("id"),
		AnswerID: data.AnswerID,
		UserID:   userIDS,
	})
	if err != nil {
		return client.schemas.PollVote.ErrorHandler(err)
	}

	return echoContext.JSON(http.StatusOK, Response{
		Success: true,
		Data:    client.schemas.PollVote.Response(resp),
		Error:   nil,
	})
}

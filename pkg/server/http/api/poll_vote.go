package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

type PollVoteModel interface {
	Do(data *poll.VoteRequest) (*poll.VoteResponse, error)
}

type PollVoteSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(*poll.VoteResponse) *schema.PollVoteResponse
}

func (client *API) PollVote(echoContext echo.Context) error {
	userID := echoContext.Get("userID")

	userIDS, ok := userID.(string)
	if !ok {
		return ErrUserIDIsNotString.WithInternal(fmt.Errorf("got %#v for userID", userID))
	}

	var data schema.PollVoteRequest

	err := json.NewDecoder(echoContext.Request().Body).Decode(&data)
	if err != nil {
		return schema.ErrPollVoteJSONDecode.WithInternal(err)
	}

	resp, err := client.models.PollVote.Do(&poll.VoteRequest{
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

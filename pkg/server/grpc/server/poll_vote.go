package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
)

type PollVoteSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *pollmodel.VoteResponse) *pollproto.VoteResponse
}

func (server *Server) Vote(ctx context.Context, input *pollproto.VoteRequest) (*pollproto.VoteResponse, error) {
	userID := ctx.Value(interceptor.UserID)

	userIDS, ok := userID.(string)
	if !ok {
		return nil, ErrUserIDIsNotString.WithInternal(fmt.Errorf("got %#v for userID", userID))
	}

	resp, err := server.models.PollVote.Do(&pollmodel.VoteRequest{
		PollID:   input.PollId,
		AnswerID: input.AnswerId,
		UserID:   userIDS,
	})
	if err != nil {
		return nil, server.schemas.PollVote.ErrorHandler(err)
	}

	return server.schemas.PollVote.Response(resp), nil
}

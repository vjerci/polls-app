package server

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
)

type PollDetailsSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *pollmodel.DetailsResponse) *pollproto.GetDetailsResponse
}

func (server *Server) GetDetails(ctx context.Context,
	input *pollproto.GetDetailsRequest,
) (*pollproto.GetDetailsResponse, error) {
	userID := ctx.Value(interceptor.UserID)

	userIDS, ok := userID.(string)
	if !ok {
		return nil, ErrUserIDIsNotString.WithInternal(fmt.Errorf("got %#v for userID", userID))
	}

	resp, err := server.models.PollDetails.Get(&pollmodel.DetailsRequest{
		ID:     input.Id,
		UserID: userIDS,
	})
	if err != nil {
		return nil, server.schemas.PollDetails.ErrorHandler(err)
	}

	return server.schemas.PollDetails.Response(resp), nil
}

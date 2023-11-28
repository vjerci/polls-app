package server

import (
	"context"

	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
)

type PollListSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *pollmodel.ListResponse) *pollproto.GetListResponse
}

func (server *Server) GetList(
	_ context.Context,
	input *pollproto.GetListRequest,
) (*pollproto.GetListResponse, error) {
	resp, err := server.models.PollList.Get(&pollmodel.ListRequest{
		Page: int(input.Page),
	})
	if err != nil {
		return nil, server.schemas.PollList.ErrorHandler(err)
	}

	return server.schemas.PollList.Response(resp), nil
}

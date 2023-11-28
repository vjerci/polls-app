package server

import (
	"context"

	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
)

type PollCreateSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *pollmodel.CreateResponse) *pollproto.CreateResponse
}

func (server *Server) Create(_ context.Context,
	input *pollproto.CreateRequest,
) (*pollproto.CreateResponse, error) {
	resp, err := server.models.PollCreate.Create(&pollmodel.CreateRequest{
		Name:    input.Name,
		Answers: input.Answers,
	})
	if err != nil {
		return nil, server.schemas.PollCreate.ErrorHandler(err)
	}

	return server.schemas.PollCreate.Response(resp), nil
}

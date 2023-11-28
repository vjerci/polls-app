package server

import (
	"context"

	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
)

type LoginSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *authmodel.LoginResponse) *authproto.LoginResponse
}

func (server *Server) Login(_ context.Context,
	input *authproto.LoginRequest,
) (*authproto.LoginResponse, error) {
	resp, err := server.models.Login.Do(&authmodel.LoginRequest{
		UserID: input.UserId,
	})
	if err != nil {
		return nil, server.schemas.Login.ErrorHandler(err)
	}

	return server.schemas.Login.Response(resp), nil
}

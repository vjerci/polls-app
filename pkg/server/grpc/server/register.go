package server

import (
	"context"

	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
)

type RegisterSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input auth.AccessToken) *authproto.RegisterResponse
}

func (server *Server) Register(_ context.Context,
	input *authproto.RegisterRequest,
) (*authproto.RegisterResponse, error) {
	resp, err := server.models.Register.Do(&authmodel.RegisterRequest{
		UserID: input.UserId,
		Name:   input.Name,
	})
	if err != nil {
		return nil, server.schemas.Register.ErrorHandler(err)
	}

	return server.schemas.Register.Response(resp), nil
}

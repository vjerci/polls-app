package server

import (
	"context"

	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
)

type GoogleLoginSchemaMap interface {
	ErrorHandler(err error) *echo.HTTPError
	Response(input *authmodel.GoogleLoginResponse) *authproto.GoogleLoginResponse
}

func (server *Server) GoogleLogin(_ context.Context,
	input *authproto.GoogleLoginRequest,
) (*authproto.GoogleLoginResponse, error) {
	resp, err := server.models.GoogleLogin.Do(&authmodel.GoogleLoginRequest{
		Token: input.Token,
	})
	if err != nil {
		return nil, server.schemas.GoogleLogin.ErrorHandler(err)
	}

	return server.schemas.GoogleLogin.Response(resp), nil
}

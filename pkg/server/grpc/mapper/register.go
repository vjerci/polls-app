package mapper

import (
	"github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type RegisterSchemaMap struct {
	API api.RegisterSchemaMap
}

func (mapper *RegisterSchemaMap) Response(input auth.AccessToken) *authproto.RegisterResponse {
	return &authproto.RegisterResponse{
		Token: string(input),
	}
}

func (mapper *RegisterSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

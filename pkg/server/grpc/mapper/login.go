package mapper

import (
	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type LoginSchemaMap struct {
	API api.LoginSchemaMap
}

func (mapper *LoginSchemaMap) Response(input *authmodel.LoginResponse) *authproto.LoginResponse {
	return &authproto.LoginResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

func (mapper *LoginSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

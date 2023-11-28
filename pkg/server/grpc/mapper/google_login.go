package mapper

import (
	"github.com/labstack/echo/v4"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type GoogleLoginSchemaMap struct {
	API api.GoogleLoginSchemaMap
}

func (mapper *GoogleLoginSchemaMap) Response(input *authmodel.GoogleLoginResponse) *authproto.GoogleLoginResponse {
	return &authproto.GoogleLoginResponse{
		Token: input.Token,
		Name:  input.Name,
	}
}

func (mapper *GoogleLoginSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

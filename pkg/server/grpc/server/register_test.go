package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/auth"
	authutil "github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/server"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type MockRegisterModel struct {
	InputData     *auth.RegisterRequest
	ResponseData  authutil.AccessToken
	ResponseError error
}

func (mock *MockRegisterModel) Do(input *auth.RegisterRequest) (authutil.AccessToken, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *authproto.RegisterRequest
		Model         *MockRegisterModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &authproto.RegisterRequest{
				UserId: "testUserID",
				Name:   "testName",
			},
			Model: &MockRegisterModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			Register: &mapper.RegisterSchemaMap{
				API: &schema.RegisterSchemaMap{},
			},
		}

		models := &api.Models{
			Register: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.Register(context.TODO(), test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestRegisterSuccessful(t *testing.T) {
	t.Parallel()

	input := &authproto.RegisterRequest{
		UserId: "testUserID",
		Name:   "testName",
	}

	modelMock := &MockRegisterModel{
		ResponseError: nil,
		ResponseData:  authutil.AccessToken("testToken"),
	}

	schemaMaps := &server.SchemasMap{
		Register: &mapper.RegisterSchemaMap{
			API: &schema.RegisterSchemaMap{},
		},
	}

	models := &api.Models{
		Register: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	resp, err := grpcServer.Register(context.TODO(), input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

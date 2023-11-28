package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
	authproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/server"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type MockLoginModel struct {
	InputData     *auth.LoginRequest
	ResponseData  *auth.LoginResponse
	ResponseError error
}

func (mock *MockLoginModel) Do(input *auth.LoginRequest) (*auth.LoginResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *authproto.LoginRequest
		Model         *MockLoginModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &authproto.LoginRequest{
				UserId: "testUserID",
			},
			Model: &MockLoginModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			Login: &mapper.LoginSchemaMap{
				API: &schema.LoginSchemaMap{},
			},
		}

		models := &api.Models{
			Login: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.Login(context.TODO(), test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestLoginSuccessful(t *testing.T) {
	t.Parallel()

	input := &authproto.LoginRequest{
		UserId: "testUserID",
	}

	modelMock := &MockLoginModel{
		ResponseError: nil,
		ResponseData: &auth.LoginResponse{
			Token: "testToken",
			Name:  "testName",
		},
	}

	schemaMaps := &server.SchemasMap{
		Login: &mapper.LoginSchemaMap{
			API: &schema.LoginSchemaMap{},
		},
	}

	models := &api.Models{
		Login: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	resp, err := grpcServer.Login(context.TODO(), input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

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

type MockGoogleLoginModel struct {
	InputData     *auth.GoogleLoginRequest
	ResponseData  *auth.GoogleLoginResponse
	ResponseError error
}

func (mock *MockGoogleLoginModel) Do(input *auth.GoogleLoginRequest) (*auth.GoogleLoginResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestGoogleLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *authproto.GoogleLoginRequest
		Model         *MockGoogleLoginModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &authproto.GoogleLoginRequest{
				Token: "testToken",
			},
			Model: &MockGoogleLoginModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			GoogleLogin: &mapper.GoogleLoginSchemaMap{
				API: &schema.GoogleLoginSchemaMap{},
			},
		}

		models := &api.Models{
			GoogleLogin: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.GoogleLogin(context.TODO(), test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestGoogleLoginSuccessful(t *testing.T) {
	t.Parallel()

	input := &authproto.GoogleLoginRequest{
		Token: "testToken",
	}

	modelMock := &MockGoogleLoginModel{
		ResponseError: nil,
		ResponseData: &auth.GoogleLoginResponse{
			Token: "testToken",
			Name:  "testName",
		},
	}

	schemaMaps := &server.SchemasMap{
		GoogleLogin: &mapper.GoogleLoginSchemaMap{
			API: &schema.GoogleLoginSchemaMap{},
		},
	}

	models := &api.Models{
		GoogleLogin: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	resp, err := grpcServer.GoogleLogin(context.TODO(), input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/server"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type MockPollDetailsModel struct {
	InputData     *poll.DetailsRequest
	ResponseData  *poll.DetailsResponse
	ResponseError error
}

func (mock *MockPollDetailsModel) Get(input *poll.DetailsRequest) (*poll.DetailsResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *pollproto.GetDetailsRequest
		//nolint:containedctx
		ctx   context.Context
		Model *MockPollDetailsModel
	}{
		{
			ExpectedError: server.ErrUserIDIsNotString,
			Input: &pollproto.GetDetailsRequest{
				Id: "testId",
			},
			ctx: context.TODO(),
		},
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &pollproto.GetDetailsRequest{
				Id: "testId",
			},
			ctx: context.WithValue(context.TODO(), interceptor.UserID, "test"),
			Model: &MockPollDetailsModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			PollDetails: &mapper.PollDetailsSchemaMap{
				API: &schema.PollDetailsSchemaMap{},
			},
		}

		models := &api.Models{
			PollDetails: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.GetDetails(test.ctx, test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollDetailsSuccessful(t *testing.T) {
	t.Parallel()

	input := &pollproto.GetDetailsRequest{
		Id: "testPollID",
	}

	modelMock := &MockPollDetailsModel{
		ResponseError: nil,
		ResponseData: &poll.DetailsResponse{
			ID:         "testPollID",
			Name:       "testPollName",
			UserAnswer: "testUserAnswerID",
			Answers: []poll.DetailsAnswer{
				{
					Name:       "testAnswerName",
					ID:         "testAnswerID",
					VotesCount: 2,
				},
			},
		},
	}

	schemaMaps := &server.SchemasMap{
		PollDetails: &mapper.PollDetailsSchemaMap{
			API: &schema.PollDetailsSchemaMap{},
		},
	}

	models := &api.Models{
		PollDetails: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	ctx := context.WithValue(context.TODO(), interceptor.UserID, "testUserID")

	resp, err := grpcServer.GetDetails(ctx, input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

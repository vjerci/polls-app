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

type MockPollVoteModel struct {
	InputData     *poll.VoteRequest
	ResponseData  *poll.VoteResponse
	ResponseError error
}

func (mock *MockPollVoteModel) Do(input *poll.VoteRequest) (*poll.VoteResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *pollproto.VoteRequest
		//nolint:containedctx
		ctx   context.Context
		Model *MockPollVoteModel
	}{
		{
			ExpectedError: server.ErrUserIDIsNotString,
			Input: &pollproto.VoteRequest{
				PollId:   "testID",
				AnswerId: "testAnswerID",
			},
			ctx: context.TODO(),
		},
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &pollproto.VoteRequest{
				PollId:   "testID",
				AnswerId: "testAnswerID",
			},
			ctx: context.WithValue(context.TODO(), interceptor.UserID, "test"),
			Model: &MockPollVoteModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			PollVote: &mapper.PollVoteSchemaMap{
				API: &schema.PollVoteSchemaMap{},
			},
		}

		models := &api.Models{
			PollVote: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.Vote(test.ctx, test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollVoteSuccessful(t *testing.T) {
	t.Parallel()

	input := &pollproto.VoteRequest{
		PollId:   "pollID",
		AnswerId: "answerID",
	}

	modelMock := &MockPollVoteModel{
		ResponseError: nil,
		ResponseData: &poll.VoteResponse{
			ModifiedAnswer: true,
		},
	}

	schemaMaps := &server.SchemasMap{
		PollVote: &mapper.PollVoteSchemaMap{
			API: &schema.PollVoteSchemaMap{},
		},
	}

	models := &api.Models{
		PollVote: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	ctx := context.WithValue(context.TODO(), interceptor.UserID, "testUserID")

	resp, err := grpcServer.Vote(ctx, input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

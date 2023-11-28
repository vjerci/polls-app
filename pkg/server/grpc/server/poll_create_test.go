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
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

type MockPollCreateModel struct {
	InputData     *poll.CreateRequest
	ResponseData  *poll.CreateResponse
	ResponseError error
}

func (mock *MockPollCreateModel) Create(input *poll.CreateRequest) (*poll.CreateResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *pollproto.CreateRequest
		Model         *MockPollCreateModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &pollproto.CreateRequest{
				Name:    "testName",
				Answers: []string{"testAnswer1", "testAnswer2"},
			},
			Model: &MockPollCreateModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			PollCreate: &mapper.PollCreateSchemaMap{
				API: &schema.PollCreateSchemaMap{},
			},
		}

		models := &api.Models{
			PollCreate: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.Create(context.TODO(), test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollCreateSuccessful(t *testing.T) {
	t.Parallel()

	input := &pollproto.CreateRequest{
		Name:    "testName",
		Answers: []string{"answer1", "answer2"},
	}

	modelMock := &MockPollCreateModel{
		ResponseError: nil,
		ResponseData: &poll.CreateResponse{
			PollID:     "testPollID",
			AnswersIDS: []string{"testAnswer1", "testAnswer2"},
		},
	}

	schemaMaps := &server.SchemasMap{
		PollCreate: &mapper.PollCreateSchemaMap{
			API: &schema.PollCreateSchemaMap{},
		},
	}

	models := &api.Models{
		PollCreate: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	resp, err := grpcServer.Create(context.TODO(), input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

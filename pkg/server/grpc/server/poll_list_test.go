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

type MockPollListModel struct {
	InputData     *poll.ListRequest
	ResponseData  *poll.ListResponse
	ResponseError error
}

func (mock *MockPollListModel) Get(input *poll.ListRequest) (*poll.ListResponse, error) {
	mock.InputData = input

	return mock.ResponseData, mock.ResponseError
}

func TestPollListErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		Input         *pollproto.GetListRequest
		Model         *MockPollListModel
	}{
		{
			ExpectedError: schema.ErrGoogleLoginModel,
			Input: &pollproto.GetListRequest{
				Page: 0,
			},
			Model: &MockPollListModel{
				ResponseError: errors.New("test error"),
			},
		},
	}

	for _, test := range testCases {
		schemaMaps := &server.SchemasMap{
			PollList: &mapper.PollListSchemaMap{
				API: &schema.PollListSchemaMap{},
			},
		}

		models := &api.Models{
			PollList: test.Model,
		}

		grpcServer := server.NewServer(models, schemaMaps)

		_, err := grpcServer.GetList(context.TODO(), test.Input)

		//nolint:errorlint
		errHTTP, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatal("expected http error")
		}

		assert.EqualValues(t, test.ExpectedError.Code, errHTTP.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, errHTTP.Message, "expected error message to match")
	}
}

func TestPollListSuccessful(t *testing.T) {
	t.Parallel()

	input := &pollproto.GetListRequest{
		Page: 0,
	}

	modelMock := &MockPollListModel{
		ResponseError: nil,
		ResponseData: &poll.ListResponse{
			HasNext: true,
			Polls: []poll.GeneralInfo{
				{
					Name: "test poll Name?",
					ID:   "testPollID",
				},
			},
		},
	}

	schemaMaps := &server.SchemasMap{
		PollList: &mapper.PollListSchemaMap{
			API: &schema.PollListSchemaMap{},
		},
	}

	models := &api.Models{
		PollList: modelMock,
	}

	grpcServer := server.NewServer(models, schemaMaps)

	resp, err := grpcServer.GetList(context.TODO(), input)

	if err != nil {
		t.Fatalf("expected returned error to be nil got '%#v' instead", err)
	}

	if resp == nil {
		t.Fatalf("expected to get a response got nil instead")
	}
}

package interceptor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MockLoginRepo struct {
	Input auth.AccessToken

	Response    string
	ResponseErr error
}

func (mock *MockLoginRepo) Decode(input auth.AccessToken) (userID string, err error) {
	mock.Input = input

	return mock.Response, mock.ResponseErr
}

type MockUserRepo struct {
	Input string

	Response    string
	ResponseErr error
}

func (mock *MockUserRepo) GetUser(userID string) (name string, err error) {
	mock.Input = userID

	return mock.Response, mock.ResponseErr
}

type MockHandler struct {
	//nolint:containedctx
	InputCtx      context.Context
	InputData     any
	ResponseData  any
	ResponseError error
}

func (mock *MockHandler) Handle(ctx context.Context, data any) (any, error) {
	mock.InputCtx = ctx
	mock.InputData = data

	return mock.ResponseData, mock.ResponseError
}

func TestAuth(t *testing.T) {
	t.Parallel()

	mockError := errors.New("test")

	testCases := []struct {
		Name string

		MetaData metadata.MD
		AuthMock interceptor.AuthRepo
		UserMock interceptor.UserRepo

		ExpectedError error
		ExpectedField string
	}{
		{
			Name: "Authentication header missing",

			MetaData:      nil,
			ExpectedError: schema.ErrAuthMissing,
		},
		{
			Name: "Authentication header invalid",

			MetaData: metadata.MD{
				"Authorization": []string{"test"},
			},
			ExpectedError: schema.ErrAuthInvalid,
		},
		{
			Name: "Auth repo error",

			MetaData: metadata.MD{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				ResponseErr: errors.New("test error"),
			},
			ExpectedError: schema.ErrAuthInvalid,
		},
		{
			Name: "User repo - user not found failure",

			MetaData: metadata.MD{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock: &MockUserRepo{
				ResponseErr: authmodel.ErrGetUserUserNotFound,
			},
			ExpectedError: schema.ErrAuthUserDoesNotExist,
		},
		{
			Name: "User repo - other types of failure",

			MetaData: metadata.MD{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock: &MockUserRepo{
				ResponseErr: mockError,
			},
			ExpectedError: mockError,
		},
		{
			Name: "Passing test",

			MetaData: metadata.MD{
				"Authorization": []string{"Bearer test"},
			},
			AuthMock: &MockLoginRepo{
				Response: "testUserID1",
			},
			UserMock:      &MockUserRepo{},
			ExpectedField: "testUserID1",
			ExpectedError: nil,
		},
	}

	for _, test := range testCases {
		testCase := test

		t.Run(t.Name()+"-"+testCase.Name, func(t *testing.T) {
			t.Parallel()

			interceptorClient := interceptor.Client{
				AuthRepo: testCase.AuthMock,
				UserRepo: testCase.UserMock,
			}

			ctx := metadata.NewIncomingContext(context.TODO(), testCase.MetaData)
			req := struct{}{}

			grpcInfo := &grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "/polls.Polls/testMethod",
			}

			handler := MockHandler{
				ResponseData:  struct{}{},
				ResponseError: nil,
				InputCtx:      nil,
				InputData:     nil,
			}

			_, err := interceptorClient.EnsureValidToken(ctx, req, grpcInfo, handler.Handle)

			assert.EqualValues(t, testCase.ExpectedError, err, "expected errors to match")

			if testCase.ExpectedField != "" {
				assert.EqualValues(t,
					testCase.ExpectedField,
					handler.InputCtx.Value(interceptor.UserID),
					"expected userID to be set")
			}
		})
	}
}

func TestSkippingAuthenticationInterceptor(t *testing.T) {
	t.Parallel()

	interceptorClient := interceptor.Client{}

	ctx := context.TODO()
	req := struct{}{}

	grpcInfo := &grpc.UnaryServerInfo{
		Server:     "test",
		FullMethod: "/polls.Auth/testMethod",
	}

	handler := &MockHandler{
		ResponseData:  struct{}{},
		ResponseError: nil,
		InputCtx:      nil,
		InputData:     nil,
	}

	//nolint:errcheck
	interceptorClient.EnsureValidToken(ctx, req, grpcInfo, handler.Handle)

	assert.EqualValues(t, req, handler.InputData, "expected handler input data to match")
}

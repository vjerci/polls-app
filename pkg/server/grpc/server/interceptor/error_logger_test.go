package interceptor_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
	"google.golang.org/grpc"
)

func TestErrorLogger(t *testing.T) {
	t.Parallel()

	internal := "internalMessage"

	testCases := []struct {
		ShouldLog       bool
		HandlerError    error
		ShouldNotReturn string
	}{
		{
			ShouldLog:    false,
			HandlerError: nil,
		},
		{
			ShouldLog:       true,
			ShouldNotReturn: internal,
			HandlerError: &echo.HTTPError{
				Message:  "errorMessage",
				Internal: errors.New(internal),
			},
		},
		{
			ShouldLog:    true,
			HandlerError: errors.New("testError"),
		},
	}

	for _, testCase := range testCases {
		buffer := bytes.NewBuffer(nil)

		interceptorClient := interceptor.Client{
			Logger: log.New(buffer, "", 0),
		}

		ctx := context.TODO()
		req := struct{}{}

		grpcInfo := &grpc.UnaryServerInfo{
			Server:     "test",
			FullMethod: "/polls.Auth/testMethod",
		}

		handler := &MockHandler{
			ResponseData:  struct{}{},
			ResponseError: testCase.HandlerError,
			InputCtx:      nil,
			InputData:     nil,
		}

		//nolint:errcheck
		interceptorClient.LogErrors(ctx, req, grpcInfo, handler.Handle)

		bytes, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read from log buffer")
		}

		if err != nil && strings.Contains(err.Error(), testCase.ShouldNotReturn) {
			t.Fatalf("internal information should't be returned")
		}

		if testCase.ShouldLog {
			if len(bytes) == 0 {
				t.Fatalf("expected to log but got nothing")
			}
		}

		if !testCase.ShouldLog {
			if len(bytes) != 0 {
				t.Fatalf("expected not to log but got %s", string(bytes))
			}
		}
	}
}

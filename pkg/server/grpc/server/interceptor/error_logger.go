package interceptor

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func (client *Client) LogErrors(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	//nolint:errorlint
	if httpError, ok := err.(*echo.HTTPError); ok {
		log.Default().Printf(
			`serving grpc error containing non visible details "%s":"%s"`,
			httpError.Message,
			httpError.Internal,
		)

		return resp, fmt.Errorf("%s", httpError.Message)
	}

	log.Default().Printf(`serving grpc error "%s"`, err)

	return resp, err
}

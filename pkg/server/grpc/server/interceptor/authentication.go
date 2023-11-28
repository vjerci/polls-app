package interceptor

import (
	"context"
	"errors"
	"strings"

	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type key int

const (
	UserID           key = iota
	pollsServiceName     = "/polls.Polls/"
)

func (client *Client) EnsureValidToken(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if !strings.HasPrefix(info.FullMethod, pollsServiceName) {
		return handler(ctx, req)
	}

	metadata, metadataExists := metadata.FromIncomingContext(ctx)
	if !metadataExists {
		return nil, schema.ErrAuthMissing
	}

	header, authorizationExists := metadata["authorization"]
	if !authorizationExists {
		return nil, schema.ErrAuthMissing
	}

	if len(header) != 1 {
		return nil, schema.ErrAuthInvalid
	}

	token, found := strings.CutPrefix(header[0], "Bearer ")
	if !found {
		return nil, schema.ErrAuthInvalid
	}

	userID, err := client.AuthRepo.Decode(auth.AccessToken(token))
	if err != nil {
		return nil, schema.ErrAuthInvalid
	}

	_, err = client.UserRepo.GetUser(userID)
	if err != nil {
		if errors.Is(err, authmodel.ErrGetUserUserNotFound) {
			return nil, schema.ErrAuthUserDoesNotExist
		}

		return nil, err
	}

	enrichedCtx := context.WithValue(ctx, UserID, userID)

	return handler(enrichedCtx, req)
}

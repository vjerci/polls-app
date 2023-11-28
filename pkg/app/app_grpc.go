package app

import (
	"errors"

	"github.com/vjerci/polls-app/pkg/config"
	"github.com/vjerci/polls-app/pkg/domain/db"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/googleauth"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
	"github.com/vjerci/polls-app/pkg/server/grpc/server"
	"github.com/vjerci/polls-app/pkg/server/grpc/server/interceptor"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"google.golang.org/grpc"
)

var ErrGRPCPortUndefined = errors.New("grpc port is undefined")

func NewGrpc(settings config.Config) (*grpc.Server, error) {
	if settings.GRPCPort == "" {
		return nil, ErrGRPCPortUndefined
	}

	dbClient := db.New(settings.PostgresURL)

	if err := dbClient.Connect(); err != nil {
		return nil, errors.Join(ErrDBConnect, err)
	}

	authClient := auth.New(settings.JWTSigningKey)
	googleAuthClient := googleauth.New(settings.GoogleClientID)

	models := newModel(dbClient, authClient, googleAuthClient)
	apiSchemaMaps := newSchemaMap()

	interceptorClient := interceptor.Client{
		AuthRepo: authClient,
		UserRepo: dbClient,
	}

	interceptors := []grpc.UnaryServerInterceptor{
		interceptorClient.EnsureValidToken,
		interceptorClient.LogErrors,
	}

	grpcServer := server.Build(models, newGrpcSchemaMaps(apiSchemaMaps), interceptors)

	return grpcServer, nil
}

func newGrpcSchemaMaps(input *api.SchemaMap) *server.SchemasMap {
	return &server.SchemasMap{
		Login:       &mapper.LoginSchemaMap{API: input.Login},
		GoogleLogin: &mapper.GoogleLoginSchemaMap{API: input.GoogleLogin},
		Register:    &mapper.RegisterSchemaMap{API: input.Register},

		PollList:    &mapper.PollListSchemaMap{API: input.PollList},
		PollDetails: &mapper.PollDetailsSchemaMap{API: input.PollDetails},
		PollCreate:  &mapper.PollCreateSchemaMap{API: input.PollCreate},
		PollVote:    &mapper.PollVoteSchemaMap{API: input.PollVote},
	}
}

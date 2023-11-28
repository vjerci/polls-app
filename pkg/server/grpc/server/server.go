package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/server/grpc/proto/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"google.golang.org/grpc"
)

type Server struct {
	poll.UnimplementedPollsServer
	auth.UnimplementedAuthServer

	models  *api.Models
	schemas *SchemasMap
}

type SchemasMap struct {
	Login       LoginSchemaMap
	GoogleLogin GoogleLoginSchemaMap
	Register    RegisterSchemaMap

	PollList    PollListSchemaMap
	PollDetails PollDetailsSchemaMap
	PollVote    PollVoteSchemaMap
	PollCreate  PollCreateSchemaMap
}

var ErrUserIDIsNotString = &echo.HTTPError{
	Message:  "internal server error, user_id is not string",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func Build(models *api.Models, schemasMap *SchemasMap, interceptors []grpc.UnaryServerInterceptor) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors...),
	}

	server := grpc.NewServer(opts...)

	implementation := NewServer(models, schemasMap)

	poll.RegisterPollsServer(server, implementation)
	auth.RegisterAuthServer(server, implementation)

	return server
}

func NewServer(models *api.Models, schemasMap *SchemasMap) *Server {
	return &Server{
		models:                   models,
		schemas:                  schemasMap,
		UnimplementedPollsServer: poll.UnimplementedPollsServer{},
		UnimplementedAuthServer:  auth.UnimplementedAuthServer{},
	}
}

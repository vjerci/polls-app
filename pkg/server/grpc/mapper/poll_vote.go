package mapper

import (
	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type PollVoteSchemaMap struct {
	API api.PollVoteSchemaMap
}

func (mapper *PollVoteSchemaMap) Response(input *pollmodel.VoteResponse) *pollproto.VoteResponse {
	return &pollproto.VoteResponse{
		ModifiedAnswer: input.ModifiedAnswer,
	}
}

func (mapper *PollVoteSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

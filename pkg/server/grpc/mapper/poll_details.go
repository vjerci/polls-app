package mapper

import (
	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type PollDetailsSchemaMap struct {
	API api.PollDetailsSchemaMap
}

func (mapper *PollDetailsSchemaMap) Response(input *pollmodel.DetailsResponse) *pollproto.GetDetailsResponse {
	//nolint: nosnakecase
	answers := make([]*pollproto.GetDetailsResponse_Answer, len(input.Answers))
	for i, v := range input.Answers {
		//nolint: nosnakecase
		answers[i] = &pollproto.GetDetailsResponse_Answer{
			Id:         v.ID,
			Name:       v.Name,
			VotesCount: int32(v.VotesCount),
		}
	}

	return &pollproto.GetDetailsResponse{
		Id:       input.ID,
		Name:     input.Name,
		UserVote: input.UserAnswer,
		Answers:  answers,
	}
}

func (mapper *PollDetailsSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

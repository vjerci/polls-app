package mapper

import (
	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type PollCreateSchemaMap struct {
	API api.PollCreateSchemaMap
}

func (mapper *PollCreateSchemaMap) Response(input *pollmodel.CreateResponse) *pollproto.CreateResponse {
	return &pollproto.CreateResponse{
		Id:        input.PollID,
		AnswerIds: input.AnswersIDS,
	}
}

func (mapper *PollCreateSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

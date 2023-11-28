package mapper

import (
	"github.com/labstack/echo/v4"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	pollproto "github.com/vjerci/polls-app/pkg/server/grpc/proto/poll"
	"github.com/vjerci/polls-app/pkg/server/http/api"
)

type PollListSchemaMap struct {
	API api.PollListSchemaMap
}

func (mapper *PollListSchemaMap) Response(input *pollmodel.ListResponse) *pollproto.GetListResponse {
	//nolint: nosnakecase
	polls := make([]*pollproto.GetListResponse_PollInfo, len(input.Polls))
	for i, v := range input.Polls {
		//nolint: nosnakecase
		polls[i] = &pollproto.GetListResponse_PollInfo{
			Id:   v.ID,
			Name: v.Name,
		}
	}

	return &pollproto.GetListResponse{
		Polls:       polls,
		HasNextPage: input.HasNext,
	}
}

func (mapper *PollListSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	return mapper.API.ErrorHandler(err)
}

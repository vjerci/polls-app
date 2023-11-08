package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
)

type PollListSchemaMap struct{}

type PollListResponse struct {
	Polls   []GeneralPollInfo `json:"polls"`
	HasNext bool              `json:"has_next"`
}

type GeneralPollInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (mapper *PollListSchemaMap) Response(input *poll.ListResponse) *PollListResponse {
	polls := make([]GeneralPollInfo, len(input.Polls))
	for i, v := range input.Polls {
		polls[i] = GeneralPollInfo{
			Name: v.Name,
			ID:   v.ID,
		}
	}

	return &PollListResponse{
		Polls:   polls,
		HasNext: input.HasNext,
	}
}

var ErrPollListInvalidPage = &echo.HTTPError{
	Message:  "invalid page specified",
	Code:     http.StatusBadRequest,
	Internal: nil,
}
var ErrPollListNoData = &echo.HTTPError{
	Message:  "poll list data for a given page does not exist",
	Code:     http.StatusNotFound,
	Internal: nil,
}

var ErrPollListModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

var ErrPollListPageNotInt = &echo.HTTPError{
	Message:  "specified page is not an integer",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

func (mapper *PollListSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, poll.ErrListInvalidPage) {
		return ErrPollListInvalidPage.WithInternal(err)
	}

	if errors.Is(err, poll.ErrListNoPolls) {
		return ErrPollListNoData.WithInternal(err)
	}

	return ErrPollListModel.WithInternal(err)
}

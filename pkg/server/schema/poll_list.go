package schema

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
)

type PollListRequest struct {
	Page int `json:"page"`
}

func (req *PollListRequest) ToModel() *model.PollListRequest {
	return &model.PollListRequest{
		Page: req.Page,
	}
}

type PollListSchemaMap struct{}

type PollListResponse struct {
	Polls   []GeneralPollInfo `json:"polls"`
	HasNext bool              `json:"has_next"`
}

type GeneralPollInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (mapper *PollListSchemaMap) Response(input *model.PollListResponse) *PollListResponse {
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

var ErrPollListJSONDecode = &echo.HTTPError{
	Message:  "failed to decode poll list json body",
	Code:     http.StatusBadRequest,
	Internal: nil,
}

var ErrPollListModel = &echo.HTTPError{
	Message:  "internal server error",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

func (mapper *PollListSchemaMap) ErrorHandler(err error) *echo.HTTPError {
	if errors.Is(err, model.ErrPollListInvalidPage) {
		return ErrPollListInvalidPage.WithInternal(err)
	}

	if errors.Is(err, model.ErrPollListNoPolls) {
		return ErrPollListNoData.WithInternal(err)
	}

	return ErrPollListModel.WithInternal(err)
}

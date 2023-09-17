package schema

import (
	"errors"
	"net/http"

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
	Polls []GeneralPollInfo `json:"polls"`
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
		Polls: polls,
	}
}

var ErrPollListInvalidPage = &UserVisibleError{
	Err:    model.ErrPollListInvalidPage,
	Status: http.StatusBadRequest,
}
var ErrPollListNoData = &UserVisibleError{
	Err:    model.ErrPollListNoPolls,
	Status: http.StatusNotFound,
}
var handledPollListErrors = []*UserVisibleError{
	ErrPollListInvalidPage,
	ErrPollListNoData,
}

var ErrPollListJSONDecode = &UserVisibleError{
	Err:    errors.New("failed to decode poll list json body"),
	Status: http.StatusBadRequest,
}

var ErrPollListModel = errors.New("model failed to get poll list")

func (mapper *PollListSchemaMap) ErrorHandler(err error) error {
	for _, targetError := range handledPollListErrors {
		if errors.Is(err, targetError.Err) {
			return targetError
		}
	}

	return errors.Join(ErrPollListModel, err)
}

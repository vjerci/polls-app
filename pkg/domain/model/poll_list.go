package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
)

type PollListRequest struct {
	Page int
}

type PollListResponse struct {
	Polls []GeneralPollInfo
}

type GeneralPollInfo struct {
	Name string
	ID   string
}

type PollListRepository interface {
	GetPollList(page int) ([]db.PollListData, error)
}

var ErrPollListInvalidPage = errors.New("invalid poll list page specified ")
var ErrPollListNoPolls = errors.New("error getting poll list, no data available")
var ErrPollListDB = errors.New("error getting poll list data")

func (client *Client) GetPollList(data *PollListRequest) (*PollListResponse, error) {
	if data.Page < 0 {
		return nil, ErrPollListInvalidPage
	}

	dbPolls, err := client.PollListDB.GetPollList(data.Page)
	if err != nil {
		return nil, errors.Join(ErrPollListDB, err)
	}

	if len(dbPolls) == 0 {
		return nil, ErrPollListNoPolls
	}

	pollInfos := make([]GeneralPollInfo, len(dbPolls))
	for i, dbPoll := range dbPolls {
		pollInfos[i] = GeneralPollInfo{
			Name: dbPoll.Name,
			ID:   dbPoll.ID,
		}
	}

	return &PollListResponse{
		Polls: pollInfos,
	}, nil
}

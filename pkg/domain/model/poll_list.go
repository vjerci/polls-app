package model

import (
	"errors"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
)

type PollListRequest struct {
	Page int
}

type PollListResponse struct {
	Polls   []GeneralPollInfo
	HasNext bool
}

type GeneralPollInfo struct {
	Name string
	ID   string
}

type PollListRepository interface {
	GetPollList(page int) ([]db.PollListData, error)
}

type PollCountRepository interface {
	HasNextPage(page int) (bool, error)
}

type PollListModel struct {
	PollListDB          PollListRepository
	PollCountRepository PollCountRepository
}

var ErrPollListInvalidPage = errors.New("invalid poll list page specified ")
var ErrPollListNoPolls = errors.New("error getting poll list, no data available")
var ErrPollListDB = errors.New("error getting poll list data")
var ErrPollListDBNextPage = errors.New("error getting next page from db")

func (model *PollListModel) Get(data *PollListRequest) (*PollListResponse, error) {
	if data.Page < 0 {
		return nil, ErrPollListInvalidPage
	}

	dbPolls, err := model.PollListDB.GetPollList(data.Page)
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

	hasNext, err := model.PollCountRepository.HasNextPage(data.Page)
	if err != nil {
		return nil, errors.Join(ErrPollListDBNextPage, err)
	}

	return &PollListResponse{
		Polls:   pollInfos,
		HasNext: hasNext,
	}, nil
}

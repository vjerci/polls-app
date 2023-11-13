package poll

import (
	"errors"

	"github.com/vjerci/polls-app/pkg/domain/db"
)

type ListRequest struct {
	Page int
}

type ListResponse struct {
	Polls   []GeneralInfo
	HasNext bool
}

type GeneralInfo struct {
	Name string
	ID   string
}

type ListRepository interface {
	GetPollList(page int) ([]db.PollListData, error)
}

type CountRepository interface {
	HasNextPage(page int) (bool, error)
}

type ListModel struct {
	ListDB  ListRepository
	CountDB CountRepository
}

var ErrListInvalidPage = errors.New("invalid poll list page specified ")
var ErrListNoPolls = errors.New("error getting poll list, no data available")
var ErrListDB = errors.New("error getting poll list data")
var ErrListDBNextPage = errors.New("error getting next page from db")

func (model *ListModel) Get(data *ListRequest) (*ListResponse, error) {
	if data.Page < 0 {
		return nil, ErrListInvalidPage
	}

	dbPolls, err := model.ListDB.GetPollList(data.Page)
	if err != nil {
		return nil, errors.Join(ErrListDB, err)
	}

	if len(dbPolls) == 0 {
		return nil, ErrListNoPolls
	}

	pollInfos := make([]GeneralInfo, len(dbPolls))
	for i, dbPoll := range dbPolls {
		pollInfos[i] = GeneralInfo{
			Name: dbPoll.Name,
			ID:   dbPoll.ID,
		}
	}

	hasNext, err := model.CountDB.HasNextPage(data.Page)
	if err != nil {
		return nil, errors.Join(ErrListDBNextPage, err)
	}

	return &ListResponse{
		Polls:   pollInfos,
		HasNext: hasNext,
	}, nil
}

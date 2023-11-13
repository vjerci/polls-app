package poll

import (
	"errors"
	"sync"

	"github.com/vjerci/polls-app/pkg/domain/db"
)

type DetailsRequest struct {
	ID     string
	UserID string
}

type DetailsResponse struct {
	ID         string
	Name       string
	UserAnswer string
	Answers    []DetailsAnswer
}

type DetailsAnswer struct {
	Name       string
	ID         string
	VotesCount int
}

type DetailsRepository interface {
	GetPollDetails(pollID string) (*db.PollDetailsResponse, error)
	GetPollDetailsAnswers(pollID string) ([]db.PollDetailsAnswer, error)
	GetUserAnswer(pollID, userID string) (answerID string, err error)
}

type DetailsModel struct {
	DetailsDB DetailsRepository
}

var ErrDetailsIDEmpty = errors.New("poll details failure, PollID cannot be empty")
var ErrDetailsUserIDEmpty = errors.New("poll details failure, UserID cannot be empty")

var ErrDetailsNoPoll = errors.New("couldn't find a poll with a given id")

var ErrDetailsQueryInfo = errors.New("failed to get poll details info")
var ErrDetailsAnswers = errors.New("failed to get answers")
var ErrDetailsUserAnswer = errors.New("failed to get user answer")

func (model *DetailsModel) Get(data *DetailsRequest) (*DetailsResponse, error) {
	if data.ID == "" {
		return nil, ErrDetailsIDEmpty
	}

	if data.UserID == "" {
		return nil, ErrDetailsUserIDEmpty
	}

	fetcher := newDetailsFetcher(data.UserID, data.ID, model.DetailsDB)

	return fetcher.Fetch()
}

type detailsFetcher struct {
	userID string
	pollID string

	db DetailsRepository

	response *DetailsResponse

	waitGroup sync.WaitGroup

	errDetailsQuery    error
	errAnswersQuery    error
	errUserAnswerQuery error
}

func newDetailsFetcher(userID string, pollID string, db DetailsRepository) *detailsFetcher {
	return &detailsFetcher{
		userID: userID,
		pollID: pollID,

		db: db,

		response: &DetailsResponse{
			Name:       "",
			UserAnswer: "",
			Answers:    nil,
			ID:         pollID,
		},

		waitGroup: sync.WaitGroup{},

		errDetailsQuery:    nil,
		errAnswersQuery:    nil,
		errUserAnswerQuery: nil,
	}
}

func (fetcher *detailsFetcher) Fetch() (*DetailsResponse, error) {
	var fetcherConcurrentRequests = 3

	fetcher.waitGroup.Add(fetcherConcurrentRequests)

	go fetcher.getPollDetails()
	go fetcher.getAnswers()
	go fetcher.getUserAnswer()

	fetcher.waitGroup.Wait()

	if err := fetcher.handleErrors(); err != nil {
		return nil, err
	}

	return fetcher.response, nil
}

func (fetcher *detailsFetcher) getPollDetails() {
	pollInfo, err := fetcher.db.GetPollDetails(fetcher.pollID)
	if err != nil {
		fetcher.errDetailsQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.Name = pollInfo.Name

	fetcher.waitGroup.Done()
}

func (fetcher *detailsFetcher) getAnswers() {
	answers, err := fetcher.db.GetPollDetailsAnswers(fetcher.pollID)
	if err != nil {
		fetcher.errAnswersQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.Answers = make([]DetailsAnswer, len(answers))
	for i, answer := range answers {
		fetcher.response.Answers[i] = DetailsAnswer{
			Name:       answer.Name,
			ID:         answer.ID,
			VotesCount: answer.Count,
		}
	}

	fetcher.waitGroup.Done()
}

func (fetcher *detailsFetcher) getUserAnswer() {
	userAnswer, err := fetcher.db.GetUserAnswer(fetcher.pollID, fetcher.userID)
	if err != nil {
		fetcher.errUserAnswerQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.UserAnswer = userAnswer

	fetcher.waitGroup.Done()
}

func (fetcher *detailsFetcher) handleErrors() error {
	details := fetcher.errDetailsQuery
	if details != nil {
		if errors.Is(details, db.ErrPollDetailsNotFound) {
			return errors.Join(ErrDetailsNoPoll, details)
		}

		return errors.Join(ErrDetailsQueryInfo, details)
	}

	answers := fetcher.errAnswersQuery
	if answers != nil {
		return errors.Join(ErrDetailsAnswers, answers)
	}

	userAnswer := fetcher.errUserAnswerQuery
	if userAnswer != nil {
		if !errors.Is(userAnswer, db.ErrUserAnswerNotFound) {
			return errors.Join(ErrDetailsUserAnswer, userAnswer)
		}
	}

	return nil
}

package model

import (
	"errors"
	"sync"

	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
)

type PollDetailsRequest struct {
	PollID string
	UserID string
}

type PollDetailsResponse struct {
	ID         string
	Name       string
	UserAnswer string
	Answers    []PollDetailsAnswer
}

type PollDetailsAnswer struct {
	Name       string
	ID         string
	VotesCount int
}

type PollDetailsRepository interface {
	GetPollDetails(pollID string) (*db.PollDetailsResponse, error)
	GetPollDetailsAnswers(pollID string) ([]db.PollDetailsAnswer, error)
	GetUserAnswer(pollID, userID string) (answerID string, err error)
}

type PollDetailsModel struct {
	PollDetailsDB PollDetailsRepository
}

var ErrPollDetailsIDEmpty = errors.New("poll details failure, PollID cannot be empty")
var ErrPollDetailsUserIDEmpty = errors.New("poll details failure, UserID cannot be empty")

var ErrPollDetailsNoPoll = errors.New("couldn't find a poll with a given id")

var ErrPollDetailsQueryInfo = errors.New("failed to get poll details info")
var ErrPollDetailsAnswers = errors.New("failed to get answers")
var ErrPollDetailsUserAnswer = errors.New("failed to get user answer")

func (model *PollDetailsModel) Get(data *PollDetailsRequest) (*PollDetailsResponse, error) {
	if data.PollID == "" {
		return nil, ErrPollDetailsIDEmpty
	}

	if data.UserID == "" {
		return nil, ErrPollDetailsUserIDEmpty
	}

	fetcher := newPollDetailsFetcher(data.UserID, data.PollID, model.PollDetailsDB)

	return fetcher.Fetch()
}

type pollDetailsFetcher struct {
	userID string
	pollID string

	db PollDetailsRepository

	response *PollDetailsResponse

	waitGroup sync.WaitGroup

	errDetailsQuery    error
	errAnswersQuery    error
	errUserAnswerQuery error
}

func newPollDetailsFetcher(userID string, pollID string, db PollDetailsRepository) *pollDetailsFetcher {
	return &pollDetailsFetcher{
		userID: userID,
		pollID: pollID,

		db: db,

		response: &PollDetailsResponse{
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

func (fetcher *pollDetailsFetcher) Fetch() (*PollDetailsResponse, error) {
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

func (fetcher *pollDetailsFetcher) getPollDetails() {
	pollInfo, err := fetcher.db.GetPollDetails(fetcher.pollID)
	if err != nil {
		fetcher.errDetailsQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.Name = pollInfo.Name

	fetcher.waitGroup.Done()
}

func (fetcher *pollDetailsFetcher) getAnswers() {
	answers, err := fetcher.db.GetPollDetailsAnswers(fetcher.pollID)
	if err != nil {
		fetcher.errAnswersQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.Answers = make([]PollDetailsAnswer, len(answers))
	for i, answer := range answers {
		fetcher.response.Answers[i] = PollDetailsAnswer{
			Name:       answer.Name,
			ID:         answer.ID,
			VotesCount: answer.Count,
		}
	}

	fetcher.waitGroup.Done()
}

func (fetcher *pollDetailsFetcher) getUserAnswer() {
	userAnswer, err := fetcher.db.GetUserAnswer(fetcher.pollID, fetcher.userID)
	if err != nil {
		fetcher.errUserAnswerQuery = err

		fetcher.waitGroup.Done()

		return
	}

	fetcher.response.UserAnswer = userAnswer

	fetcher.waitGroup.Done()
}

func (fetcher *pollDetailsFetcher) handleErrors() error {
	details := fetcher.errDetailsQuery
	if details != nil {
		if errors.Is(details, db.ErrPollDetailsNotFound) {
			return errors.Join(ErrPollDetailsNoPoll, details)
		}

		return errors.Join(ErrPollDetailsQueryInfo, details)
	}

	answers := fetcher.errAnswersQuery
	if answers != nil {
		return errors.Join(ErrPollDetailsAnswers, answers)
	}

	userAnswer := fetcher.errUserAnswerQuery
	if userAnswer != nil {
		if !errors.Is(userAnswer, db.ErrUserAnswerNotFound) {
			return errors.Join(ErrPollDetailsUserAnswer, userAnswer)
		}
	}

	return nil
}

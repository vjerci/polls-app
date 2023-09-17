package api

import (
	"errors"

	echo "github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

var ErrUserIDIsNotString = errors.New("user id is not string")

type Factory interface {
	Register(model RegisterModel, schemaMap RegisterSchemaMap) echo.HandlerFunc
	Login(model LoginModel, schemaMap LoginSchemaMap) echo.HandlerFunc

	PollList(pollListModel PollListModel, schemaMap PollListSchemaMap)
	PollDetails(pollDetailsModel PollDetailsModel, schemaMap PollDetailsSchemaMap)

	PollCreate(pollCreateModel PollCreateModel, schemaMap PollCreateSchemaMap)
	PollVote(pollVoteModel PollVoteModel, schemaMap PollVoteSchemaMap)
}

type FactoryImplementation struct{}

func New() *FactoryImplementation {
	return &FactoryImplementation{}
}

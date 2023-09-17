package api

import (
	"errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

var ErrUserIDIsNotString = errors.New("user id is not string")

type Client struct {
	models  *Models
	schemas *SchemaMap
}

type Models struct {
	Login    LoginModel
	Register RegisterModel

	PollList    PollListModel
	PollDetails PollDetailsModel

	PollCreate PollCreateModel
	PollVote   PollVoteModel
}

type SchemaMap struct {
	Login    LoginSchemaMap
	Register RegisterSchemaMap

	PollList    PollListSchemaMap
	PollDetails PollDetailsSchemaMap

	PollCreate PollCreateSchemaMap
	PollVote   PollVoteSchemaMap
}

func New(models *Models, schemas *SchemaMap) *Client {
	return &Client{
		models:  models,
		schemas: schemas,
	}
}

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

var ErrUserIDIsNotString = &echo.HTTPError{
	Message:  "internal server error, user_id is not string",
	Code:     http.StatusInternalServerError,
	Internal: nil,
}

type API struct {
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

func New(models *Models, schemas *SchemaMap) *API {
	return &API{
		models:  models,
		schemas: schemas,
	}
}

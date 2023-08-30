package api

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

// ErrorUserVisible is used for showing errors to users for easier api usage, they also enhance error messages with status code
type ErrorUserVisible struct {
	Err    error
	Status int
}

func (e *ErrorUserVisible) Error() string {
	return e.Err.Error()
}

func (e *ErrorUserVisible) Is(target error) bool {
	targetErrorUserVisible, ok := target.(*ErrorUserVisible)
	if ok {
		return errors.Is(e.Err, target) && targetErrorUserVisible.Status == e.Status
	}

	return errors.Is(e.Err, target)
}

func (e ErrorUserVisible) Unwrap() error {
	return e.Err
}

type Factory interface {
	Register(repo RegisterRepository) echo.HandlerFunc
}

type FactoryImplementation struct{}

func New() *FactoryImplementation {
	return &FactoryImplementation{}
}

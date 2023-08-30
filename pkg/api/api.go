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

type ErrorClientVisible struct {
	Err    error
	Status int
}

func (e *ErrorClientVisible) Error() string {
	return e.Err.Error()
}

func (e *ErrorClientVisible) Is(target error) bool {
	targetErrorClientVisible, ok := target.(*ErrorClientVisible)
	if ok {
		return errors.Is(e.Err, target) && targetErrorClientVisible.Status == e.Status
	}

	return errors.Is(e.Err, target)
}

type Factory interface {
	Register(repo RegisterRepository) echo.HandlerFunc
}

type FactoryImplementation struct{}

func New() *FactoryImplementation {
	return &FactoryImplementation{}
}

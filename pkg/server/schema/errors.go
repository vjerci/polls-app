package schema

import (
	"errors"
)

type Map struct{}

func NewSchemaMap() *Map {
	return &Map{}
}

// UserVisibleError is used for showing usage errors to users for easier api usage.
// They also enhance error messages with status code.
type UserVisibleError struct {
	Err    error
	Status int
}

func (e *UserVisibleError) Error() string {
	return e.Err.Error()
}

func (e *UserVisibleError) Is(target error) bool {
	targetErrorUserVisible, ok := target.(*UserVisibleError)
	if ok {
		return errors.Is(e.Err, target) && targetErrorUserVisible.Status == e.Status
	}

	return errors.Is(e.Err, target)
}

func (e UserVisibleError) Unwrap() error {
	return e.Err
}

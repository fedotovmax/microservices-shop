package adapter

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("already exists")

var ErrTimeout = errors.New("timeout expired")

var ErrUnavailable = errors.New("unavailable")

var ErrInternal = errors.New("internal error")

type ErrPartialUpdate struct {
	Expected int64
	Actual   int64
}

func (e *ErrPartialUpdate) Error() string {
	return fmt.Sprintf("partial update: expected %d entities, but updated %d", e.Expected, e.Actual)
}

func ParsePartial(affected, expected int64) error {
	if affected != expected {
		return &ErrPartialUpdate{
			Expected: expected,
			Actual:   affected,
		}
	}
	return nil
}

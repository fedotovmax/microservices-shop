package adapter

import (
	"errors"
)

var ErrNotFound = errors.New("entity not found")

var ErrAlreadyExists = errors.New("entity already exists")

var ErrTimeout = errors.New("timeout expired")

var ErrUnavailable = errors.New("resource unavailable")

var ErrInternal = errors.New("internal error")

var ErrNoFieldsToUpdate = errors.New("no fields to update")

var ErrUnsupported = errors.New("unsupported")

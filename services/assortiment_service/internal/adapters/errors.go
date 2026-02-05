package adapters

import "errors"

var ErrInternal = errors.New("internal error")

var ErrNotFound = errors.New("entity not found")

var ErrAlreadyExists = errors.New("entity already exists")

var ErrNoFieldsToUpdate = errors.New("no fields to update")

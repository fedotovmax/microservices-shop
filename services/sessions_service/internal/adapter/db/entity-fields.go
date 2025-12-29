package db

import (
	"errors"
	"fmt"
)

type SessionEntityFields string

func (se SessionEntityFields) String() string {
	return string(se)
}

const (
	SessionFieldID          SessionEntityFields = "id"
	SessionFieldRefreshHash SessionEntityFields = "refresh_hash"
)

var ErrSessionEntityField = errors.New("the passed field does not belong to the user entity")

func IsSessionEntityField(f SessionEntityFields) error {

	const op = "db.IsSessionEntityField"

	switch f {
	case SessionFieldID, SessionFieldRefreshHash:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrSessionEntityField)
}

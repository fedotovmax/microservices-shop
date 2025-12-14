package db

import (
	"errors"
	"fmt"
)

type UserEntityFields string

func (ue UserEntityFields) String() string {
	return string(ue)
}

const (
	UserFieldID    UserEntityFields = "id"
	UserFieldEmail UserEntityFields = "email"
)

var ErrUserEntityField = errors.New("the passed field does not belong to the user entity")

func IsUserEntityField(f UserEntityFields) error {

	const op = "db.IsUserEntityField"

	switch f {
	case UserFieldEmail, UserFieldID:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrUserEntityField)
}

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

type SecurityTable string

func (st SecurityTable) String() string {
	return string(st)
}

const (
	SecurityTableBlacklist SecurityTable = "blacklist"
	SecurityTableBypass    SecurityTable = "bypass"
)

var ErrInvalidSecurityTableName = errors.New("the passed name does not belong to security tables")

func IsSecurityTable(t SecurityTable) error {

	const op = "db.IsSecurityTable"

	switch t {
	case SecurityTableBlacklist, SecurityTableBypass:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrInvalidSecurityTableName)
}

type Operation uint8

const (
	OperationSelect Operation = iota
	OperationInsert
	OperationUpdate
	OpearionDelete
)

var ErrInvalidOperation = errors.New("invalid operation")

func IsOperation(o Operation) error {

	const op = "db.IsOperation"

	switch o {
	case OperationSelect, OperationInsert, OperationUpdate, OpearionDelete:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrInvalidOperation)
}

package errs

import "errors"

var ErrAppNotFound = errors.New("app not found")

var ErrInvalidAppType = errors.New("invalid app type")

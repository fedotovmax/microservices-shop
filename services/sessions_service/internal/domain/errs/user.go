package errs

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInternalCreateUser = errors.New("internal error when create user")

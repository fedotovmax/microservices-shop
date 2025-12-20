package errs

import "errors"

var ErrInvalidCommand = errors.New("invalid command")

var ErrUserIDAlreadyExists = errors.New("user id already exists")

var ErrChatIDAlreadyExists = errors.New("chat id already exists")

var ErrChatIDNotFound = errors.New("chat id not found")
var ErrUserIDNotFound = errors.New("user id not found")

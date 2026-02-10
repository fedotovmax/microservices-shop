package errs

import "errors"

var ErrInvalidCommand = errors.New("invalid command")

var ErrUserIDAlreadyExists = errors.New("user id already exists")

var ErrChatIDAlreadyExists = errors.New("chat id already exists")

var ErrChatIDNotFound = errors.New("chat id not found")
var ErrUserIDNotFound = errors.New("user id not found")

var ErrEventNotFound = errors.New("event not found")

var ErrEventAlreadyHandled = errors.New("event already handled")

var ErrSendTelegramMessage = errors.New("error when send message to telegram")

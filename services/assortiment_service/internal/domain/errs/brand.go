package errs

import "errors"

var ErrBrandNotFound = errors.New("brand not found")
var ErrBrandAlreadyExists = errors.New("brand already exists")

var ErrLanguageNotFound = errors.New("language not found")

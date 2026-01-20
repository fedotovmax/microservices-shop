package config

import "errors"

var ErrRequiredConfigPath = errors.New("config path is required for current APP_ENV")

var ErrConfigPathNotExists = errors.New("config path is not exists, check path or create file")

var ErrRequiredNewAppName = errors.New("name is required arg")
var ErrRequiredNewAppType = errors.New("type is required arg")

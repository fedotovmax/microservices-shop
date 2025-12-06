package config

import "errors"

var ErrRequiredConfigPath = errors.New("config path is required for current APP_ENV")

var ErrConfigPathNotExists = errors.New("config path is not exists, check path or create file")

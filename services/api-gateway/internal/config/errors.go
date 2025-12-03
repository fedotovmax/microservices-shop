package config

import "errors"

var errRequiredConfigPath = errors.New("config path is required for current APP_ENV")

var errConfigPathNotExists = errors.New("config path is not exists, check path or create file")

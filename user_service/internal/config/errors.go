package config

import "errors"

var errVariableParse = errors.New("error when parse variable from env")

var errVariableNotProvided = errors.New("env variable not provided")

var errVariableIsEmpty = errors.New("enb variable is empty")

var errInvalidAppEnv = errors.New("passed argument value \"APP_ENV\" is not supported")

var errRequiredConfigPath = errors.New("config path is required for current APP_ENV")

var errConfigPathNotExists = errors.New("config path is not exists, check path or create file")

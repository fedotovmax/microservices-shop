package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkConfigPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func parseEnv(env string) bool {
	switch env {
	case Development, Production, Local:
		return true
	default:
		return false
	}
}

func getCurrentAppEnv() string {
	const op = "config.getCurrentAppEnv"

	env, err := getEnv("APP_ENV")

	if err != nil {
		panicError(op, err)
	}

	ok := parseEnv(env)

	if !ok {
		panicError(op, errInvalidAppEnv)
	}

	return env
}

func panicError(op string, err error) string {
	return fmt.Sprintf("%s: %s", op, err.Error())
}

func getEnv(name string) (string, error) {
	const op = "config.getEnv"
	value, exists := os.LookupEnv(name)

	if !exists {
		return "", fmt.Errorf("%s: variable key: %s: %w", op, name, errVariableNotProvided)
	}

	if value == "" {
		return "", fmt.Errorf("%s: variable key: %s: %w", op, name, errVariableIsEmpty)
	}

	return value, nil
}

func getEnvAsInt(name string) (int, error) {
	const op = "config.getEnvAsInt"

	valueStr, err := getEnv(name)

	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		return 0, fmt.Errorf("%s: variable key: %s: %w: %v", op, name, errVariableParse, err)
	}

	return value, nil
}

func getEnvAsBool(name string) (bool, error) {
	const op = "config.getEnvAsBool"

	valueStr, err := getEnv(name)

	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(valueStr)

	if err != nil {
		return false, fmt.Errorf("%s: variable key: %s: %w: %v", op, name, errVariableParse, err)
	}

	return value, nil
}

func getEnvAsArr(name string) ([]string, error) {
	const op = "config.getEnvAsArr"

	valueStr, err := getEnv(name)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	arr := strings.Split(valueStr, ",")

	return arr, nil
}

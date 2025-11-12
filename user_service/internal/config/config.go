package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Port         int
	DBUrl        string
	KafkaBrokers []string
}

var errVariableParse = errors.New("error when parse variable from env")

var errVariableNotProvided = errors.New("env variable not provided")

var errVariableIsEmpty = errors.New("enb variable is empty")

var config *Config

var initErr error

var once sync.Once

func New() (*Config, error) {
	const op = "config.New"
	once.Do(func() {
		port, err := getEnvAsInt("PORT")

		if err != nil {
			initErr = fmt.Errorf("%s: %w", op, err)
			return
		}

		dbUrl, err := getEnv("DB_URL")

		if err != nil {
			initErr = fmt.Errorf("%s: %w", op, err)
			return
		}

		kafkaBrokers, err := getEnvAsArr("KAFKA_BROKERS")

		if err != nil {
			initErr = fmt.Errorf("%s: %w", op, err)
			return
		}

		config = &Config{
			Port:         port,
			DBUrl:        dbUrl,
			KafkaBrokers: kafkaBrokers,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return config, nil

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

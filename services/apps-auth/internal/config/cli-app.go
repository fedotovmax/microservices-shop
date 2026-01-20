package config

import (
	"errors"
	"fmt"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/keys"
	"github.com/fedotovmax/validation"
	"github.com/joho/godotenv"
)

type CliAuthAppConfig struct {
	RedisAddr     string
	RedisPassword string
	NewAppName    string
	NewAppType    int
}

// Load config from file, when required APP_ENV variable provided and equal to local or development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func LoadCliAuthAppConfig() (*CliAuthAppConfig, error) {
	const op = "config.LoadCliAuthAppConfig"

	appEnv, err := envconfig.GetCurrentAppEnv(keys.AppEnv, keys.SupportedEnv)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var isDev = appEnv == keys.Local || appEnv == keys.Development

	appflags, err := loadCliAppConfigFlags(isDev)

	if err != nil {
		return nil, err
	}

	if isDev {
		ok := envconfig.CheckConfigPathExists(appflags.ConfigPath)

		if !ok {
			return nil, fmt.Errorf("%s: %w", op, ErrConfigPathNotExists)
		}

		err = godotenv.Load(appflags.ConfigPath)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	redisPassword, err := envconfig.GetEnv("REDIS_PASSWORD")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	redisAddr, err := envconfig.GetEnv("REDIS_ADDR")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config := &CliAuthAppConfig{
		NewAppName:    appflags.NewAppName,
		NewAppType:    appflags.NewAppType,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
	}

	err = config.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: invalid config: %w", op, err)
	}

	return config, nil
}

func (c *CliAuthAppConfig) validate() error {
	var verrs []error

	_, err := validation.IsURI(c.RedisAddr)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "RedisAddr", err))
	}

	err = validation.MinLength(c.RedisPassword, 3)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "RedisPassword", err))
	}

	return errors.Join(verrs...)
}

package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/keys"
	"github.com/fedotovmax/validation"
	"github.com/joho/godotenv"
)

type AuthAppConfig struct {
	Env              string
	RedisAddr        string
	AdminSecret      string
	RedisPassword    string
	TokenSecret      string
	Issuer           string
	TokenExpDuration time.Duration
	Port             uint16
}

// Load config from file, when required APP_ENV variable provided and equal to local or development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func LoadAuthAppConfig() (*AuthAppConfig, error) {
	const op = "config.LoadAuthAppConfig"

	appEnv, err := envconfig.GetCurrentAppEnv(keys.AppEnv, keys.SupportedEnv)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if appEnv == keys.Local || appEnv == keys.Development {

		appflags, err := loadAppConfigFlags()

		if err != nil {
			return nil, err
		}

		ok := envconfig.CheckConfigPathExists(appflags.ConfigPath)

		if !ok {
			return nil, fmt.Errorf("%s: %w", op, ErrConfigPathNotExists)
		}

		err = godotenv.Load(appflags.ConfigPath)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	port, err := envconfig.GetEnvAs[uint16]("APPS_AUTH_APP_PORT")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	redisPassword, err := envconfig.GetEnv("REDIS_PASSWORD")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	redisAddr, err := envconfig.GetEnv("REDIS_ADDR")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	adminSecret, err := envconfig.GetEnv("SECRET_ADMIN_KEY")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokenSecret, err := envconfig.GetEnv("TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokenExpDuration, err := envconfig.GetEnvAs[time.Duration]("TOKEN_EXPIRES_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	issuer, err := envconfig.GetEnv("APPS_TOKEN_ISSUER")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config := &AuthAppConfig{
		Env:              appEnv,
		Port:             port,
		RedisAddr:        redisAddr,
		RedisPassword:    redisPassword,
		AdminSecret:      adminSecret,
		Issuer:           issuer,
		TokenSecret:      tokenSecret,
		TokenExpDuration: tokenExpDuration,
	}

	err = config.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: invalid config: %w", op, err)
	}

	return config, nil
}

func (c *AuthAppConfig) validate() error {
	var verrs []error

	err := validation.Range(c.Port, 1024, 65535)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "Port", err))
	}

	_, err = validation.IsURI(c.RedisAddr)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "RedisAddr", err))
	}

	err = validation.MinLength(c.RedisPassword, 3)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "RedisPassword", err))
	}

	err = validation.MinLength(c.Issuer, 3)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "Issuer", err))
	}

	err = validation.MinLength(c.AdminSecret, 3)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "AdminSecret", err))
	}

	err = validation.MinLength(c.TokenSecret, 3)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "TokenSecret", err))
	}

	err = validation.Min(c.TokenExpDuration, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "TokenExpDuration", err))
	}

	return errors.Join(verrs...)
}

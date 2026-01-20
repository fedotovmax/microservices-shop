package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/validation"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env                     string
	UsersClientAddr         string
	SessionsClientAddr      string
	TranslationPath         string
	SessionsTokenIssuer     string
	ApplicationsTokenIssuer string
	SessionsTokenSecret     string
	ApplicationsTokenSecret string
	Port                    uint16
}

type appFlags struct {
	ConfigPath string
}

// Load config from file, when required APP_ENV variable provided and equal to local or development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func LoadAppConfig() (*AppConfig, error) {
	const op = "config.MustLoadAppConfig"

	appEnv, err := envconfig.GetCurrentAppEnv(keys.AppEnv, keys.SupportedEnv)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if appEnv == keys.Development {

		appflags, err := loadAppConfigFlags()

		if err != nil {
			return nil, err
		}

		ok := envconfig.CheckConfigPathExists(appflags.ConfigPath)

		if !ok {
			return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
		}

		err = godotenv.Load(appflags.ConfigPath)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
		}
	}

	port, err := envconfig.GetEnvAs[uint16]("PORT")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	usersClientAddr, err := envconfig.GetEnv("USERS_CLIENT_ADDR")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	sessionsClientAddr, err := envconfig.GetEnv("SESSIONS_CLIENT_ADDR")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	translationPath, err := envconfig.GetEnv("TRANSLATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	sessionsTokenIssuer, err := envconfig.GetEnv("SESSIONS_TOKEN_ISSUER")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	sessionsTokenSecret, err := envconfig.GetEnv("SESSIONS_TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	appsTokenSecret, err := envconfig.GetEnv("APPS_TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	appsTokenIssuer, err := envconfig.GetEnv("APPS_TOKEN_ISSUER")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	config := &AppConfig{
		Env:                     appEnv,
		Port:                    port,
		UsersClientAddr:         usersClientAddr,
		SessionsClientAddr:      sessionsClientAddr,
		TranslationPath:         translationPath,
		SessionsTokenIssuer:     sessionsTokenIssuer,
		ApplicationsTokenIssuer: appsTokenIssuer,
		SessionsTokenSecret:     sessionsTokenSecret,
		ApplicationsTokenSecret: appsTokenSecret,
	}

	err = config.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: invalid config: %w", op, err)
	}

	return config, nil
}

func (c *AppConfig) validate() error {
	var verrs []error

	// TRANSLATIONS_PATH
	err := validation.IsFilePath(c.TranslationPath)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "TranslationPath", err))
	}

	// USERS_CLIENT_ADDR
	_, err = validation.IsURI(c.UsersClientAddr)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "UsersClientAddr", err))
	}

	// SESSIONS_CLIENT_ADDR
	_, err = validation.IsURI(c.SessionsClientAddr)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "SessionsClientAddr", err))
	}

	// SESSIONS_TOKEN_ISSUER
	err = validation.MinLength(c.SessionsTokenIssuer, 1)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "SessionsTokenIssuer", err))
	}

	// APPLICATIONS_TOKEN_ISSUER
	err = validation.MinLength(c.ApplicationsTokenIssuer, 1)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "ApplicationsTokenIssuer", err))
	}

	// SESSIONS_TOKEN_SECRET
	err = validation.MinLength(c.SessionsTokenSecret, 1)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "SessionsTokenSecret", err))
	}

	// APPLICATIONS_TOKEN_SECRET
	err = validation.MinLength(c.ApplicationsTokenSecret, 1)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "ApplicationsTokenSecret", err))
	}

	// PORT
	err = validation.Range(c.Port, 1024, 65535)
	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "Port", err))
	}

	return errors.Join(verrs...)
}

func loadAppConfigFlags() (*appFlags, error) {

	const op = "config.loadConfigPath"

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configPath == "" {
		return nil, fmt.Errorf("%s: %w", op, errRequiredConfigPath)
	}

	return &appFlags{
		ConfigPath: configPath,
	}, nil
}

package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	KafkaBrokers            []string
	Env                     string
	DBUrl                   string
	TranslationPath         string
	Port                    uint16
	AccessTokenSecret       string
	RefreshTokenSecret      string
	AccessTokenExpDuration  time.Duration
	RefreshTokenExpDuration time.Duration
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

	port, err := envconfig.GetEnvAs[uint16]("PORT")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dbUrl, err := envconfig.GetEnv("DB_URL")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	kafkaBrokers, err := envconfig.GetEnvAsArr[string]("KAFKA_BROKERS")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	translationPath, err := envconfig.GetEnv("TRANSLATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenSecret, err := envconfig.GetEnv("ACCESS_TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenSecret, err := envconfig.GetEnv("REFRESH_TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenExpDuration, err := envconfig.GetEnvAs[time.Duration]("ACCESS_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	refreshTokenExpDuration, err := envconfig.GetEnvAs[time.Duration]("REFRESH_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config := &AppConfig{
		Env:                     appEnv,
		Port:                    port,
		DBUrl:                   dbUrl,
		KafkaBrokers:            kafkaBrokers,
		TranslationPath:         translationPath,
		AccessTokenSecret:       accessTokenSecret,
		RefreshTokenSecret:      refreshTokenSecret,
		AccessTokenExpDuration:  accessTokenExpDuration,
		RefreshTokenExpDuration: refreshTokenExpDuration,
	}

	return config, nil
}

func loadAppConfigFlags() (*appFlags, error) {

	const op = "config.loadConfigPath"

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configPath == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredConfigPath)
	}

	return &appFlags{
		ConfigPath: configPath,
	}, nil
}

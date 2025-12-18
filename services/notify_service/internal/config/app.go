package config

import (
	"flag"
	"fmt"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	KafkaBrokers    []string
	Env             string
	TranslationPath string
	TgBotToken      string
	RedisAddr       string
	RedisPassword   string
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

	kafkaBrokers, err := envconfig.GetEnvAsArr[string]("KAFKA_BROKERS")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	translationPath, err := envconfig.GetEnv("TRANSLATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	tgBotToken, err := envconfig.GetEnv("TG_BOT_TOKEN")

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

	config := &AppConfig{
		Env:             appEnv,
		KafkaBrokers:    kafkaBrokers,
		TranslationPath: translationPath,
		TgBotToken:      tgBotToken,
		RedisAddr:       redisAddr,
		RedisPassword:   redisPassword,
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

package config

import (
	"flag"
	"fmt"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env             string
	Port            uint16
	UsersClientAddr string
	TranslationPath string
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

	translationPath, err := envconfig.GetEnv("TRANSLATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, errConfigPathNotExists)
	}

	config := &AppConfig{
		Env:             appEnv,
		Port:            port,
		UsersClientAddr: usersClientAddr,
		TranslationPath: translationPath,
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
		return nil, fmt.Errorf("%s: %w", op, errRequiredConfigPath)
	}

	return &appFlags{
		ConfigPath: configPath,
	}, nil
}

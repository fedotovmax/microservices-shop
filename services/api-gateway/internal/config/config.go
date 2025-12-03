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
func MustLoadAppConfig() *AppConfig {
	const op = "config.MustLoadAppConfig"

	appEnv, err := envconfig.GetCurrentAppEnv(keys.AppEnv, keys.SupportedEnv)

	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	if appEnv == keys.Development {

		appflags := loadAppConfigFlags()

		ok := envconfig.CheckConfigPathExists(appflags.ConfigPath)

		if !ok {
			panic(fmt.Errorf("%s: %w", op, errConfigPathNotExists))
		}

		err := godotenv.Load(appflags.ConfigPath)

		if err != nil {
			panic(fmt.Errorf("%s: %w", op, err))
		}

	}

	port, err := envconfig.GetEnvAs[uint16]("PORT")

	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	usersClientAddr, err := envconfig.GetEnv("USERS_CLIENT_ADDR")

	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	translationPath, err := envconfig.GetEnv("TRANSLATIONS_PATH")

	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	config := &AppConfig{
		Env:             appEnv,
		Port:            port,
		UsersClientAddr: usersClientAddr,
		TranslationPath: translationPath,
	}

	return config
}

func loadAppConfigFlags() *appFlags {

	const op = "config.loadConfigPath"

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configPath == "" {
		panic(fmt.Errorf("%s: %w", op, errRequiredConfigPath))
	}

	return &appFlags{
		ConfigPath: configPath,
	}
}

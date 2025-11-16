package config

import (
	"flag"

	"github.com/joho/godotenv"
)

const (
	Local       = "local"
	Development = "development"
	Production  = "production"
)

type AppConfig struct {
	Env          string
	Port         int
	DBUrl        string
	KafkaBrokers []string
}

type appFlags struct {
	ConfigPath string
}

// Load config from file, when required APP_ENV variable provided and equal to local or development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func MustLoadAppConfig() *AppConfig {
	const op = "config.MustLoadAppConfig"

	appEnv := getCurrentAppEnv()

	if appEnv == Local || appEnv == Development {

		appflags := loadAppConfigFlags()

		ok := checkConfigPathExists(appflags.ConfigPath)

		if !ok {
			panicError(op, errConfigPathNotExists)
		}

		err := godotenv.Load(appflags.ConfigPath)

		if err != nil {
			panicError(op, err)
		}

	}

	port, err := getEnvAsInt("PORT")

	if err != nil {
		panic(panicError(op, err))
	}

	dbUrl, err := getEnv("DB_URL")

	if err != nil {
		panic(panicError(op, err))
	}

	kafkaBrokers, err := getEnvAsArr("KAFKA_BROKERS")

	if err != nil {
		panic(panicError(op, err))
	}

	config := &AppConfig{
		Env:          appEnv,
		Port:         port,
		DBUrl:        dbUrl,
		KafkaBrokers: kafkaBrokers,
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
		panicError(op, errRequiredConfigPath)
	}

	return &appFlags{
		ConfigPath: configPath,
	}
}

package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/validation"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	KafkaBrokers             []string
	Env                      string
	DBUrl                    string
	TranslationPath          string
	AccessTokenSecret        string
	TokenIssuer              string
	AccessTokenExpDuration   time.Duration
	RefreshTokenExpDuration  time.Duration
	BlacklistCodeExpDuration time.Duration
	LoginBypassExpDuration   time.Duration
	Port                     uint16
	BlacklistCodeLength      uint8
	LoginBypassCodeLength    uint8
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

	accessTokenExpDuration, err := envconfig.GetEnvAs[time.Duration]("ACCESS_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	refreshTokenExpDuration, err := envconfig.GetEnvAs[time.Duration]("REFRESH_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	blacklistCodeExpDuration, err := envconfig.GetEnvAs[time.Duration]("BLACKLIST_CODE_EXP_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	loginBypassCodeExpDuration, err := envconfig.GetEnvAs[time.Duration]("LOGIN_BYPASS_CODE_EXP_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	loginBypassCodeLength, err := envconfig.GetEnvAs[uint8]("LOGIN_BYPASS_CODE_LENGTH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	blacklistCodeLength, err := envconfig.GetEnvAs[uint8]("BLACKLIST_CODE_LENGTH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokenIssuer, err := envconfig.GetEnv("TOKEN_ISSUER")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config := &AppConfig{
		Env:                      appEnv,
		Port:                     port,
		DBUrl:                    dbUrl,
		KafkaBrokers:             kafkaBrokers,
		TranslationPath:          translationPath,
		AccessTokenSecret:        accessTokenSecret,
		AccessTokenExpDuration:   accessTokenExpDuration,
		RefreshTokenExpDuration:  refreshTokenExpDuration,
		BlacklistCodeLength:      blacklistCodeLength,
		BlacklistCodeExpDuration: blacklistCodeExpDuration,
		LoginBypassCodeLength:    loginBypassCodeLength,
		LoginBypassExpDuration:   loginBypassCodeExpDuration,
		TokenIssuer:              tokenIssuer,
	}

	err = config.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: invalid config: %w", op, err)
	}

	return config, nil
}

func (c *AppConfig) validate() error {

	var verrs []error

	err := validation.IsFilePath(c.TranslationPath)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "TranslationPath", err))
	}

	err = validation.Min(c.AccessTokenExpDuration, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "AccessTokenExpDuration", err))
	}

	err = validation.Min(c.RefreshTokenExpDuration, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "RefreshTokenExpDuration", err))
	}

	err = validation.Min(c.BlacklistCodeExpDuration, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "BlacklistCodeExpDuration", err))
	}

	err = validation.Min(c.LoginBypassExpDuration, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "LoginBypassExpDuration", err))
	}

	err = validation.Min(c.BlacklistCodeLength, 6)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "BlacklistCodeLength", err))
	}

	err = validation.Min(c.LoginBypassCodeLength, 6)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "LoginBypassCodeLength", err))
	}

	_, err = validation.IsURI(c.DBUrl)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "DBUrl", err))
	}

	for idx, kafkaBroker := range c.KafkaBrokers {
		_, err = validation.IsURI(kafkaBroker)
		if err != nil {
			verrs = append(verrs, fmt.Errorf("%s[%d]: %w", "KafkaBrokers", idx, err))
		}
	}

	err = validation.MinLength(c.TokenIssuer, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "TokenIssuer", err))
	}

	err = validation.MinLength(c.AccessTokenSecret, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "AccessTokenSecret", err))
	}

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
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredConfigPath)
	}

	return &appFlags{
		ConfigPath: configPath,
	}, nil
}

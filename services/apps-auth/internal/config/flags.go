package config

import (
	"flag"
	"fmt"
)

type cliAppFlags struct {
	ConfigPath string
	NewAppName string
	NewAppType int
}

func loadCliAppConfigFlags(isDev bool) (*cliAppFlags, error) {

	const op = "config.loadCliAppConfigFlags"

	var configPath string
	var newAppName string
	var newAppType int

	if isDev {
		flag.StringVar(&configPath, "config_path", "", "path to config file")
		flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")
	}

	flag.StringVar(&newAppName, "name", "", "new app name")
	flag.IntVar(&newAppType, "type", 0, "new app type")

	flag.Parse()

	if isDev {
		if configPath == "" {
			return nil, fmt.Errorf("%s: %w", op, ErrRequiredConfigPath)
		}
	}

	if newAppName == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredNewAppName)
	}

	if newAppType == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredNewAppType)
	}

	return &cliAppFlags{
		ConfigPath: configPath,
		NewAppName: newAppName,
		NewAppType: newAppType,
	}, nil
}

type appFlags struct {
	ConfigPath string
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

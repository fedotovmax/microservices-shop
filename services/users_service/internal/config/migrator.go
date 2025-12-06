package config

import (
	"flag"
	"fmt"

	"github.com/fedotovmax/envconfig"
	"github.com/joho/godotenv"
)

type MigratorConfig struct {
	DBUrl           string
	MigrationsPath  string
	MigrationsTable string
	Cmd             *string
	Version         *int
	Steps           *int
}

func LoadMigratorConfig() (*MigratorConfig, error) {

	const op = "config.MustLoadMigratorConfig"

	mflags, err := loadMigratorFlags()

	if err != nil {
		return nil, err
	}

	ok := envconfig.CheckConfigPathExists(mflags.ConfigPath)

	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrConfigPathNotExists)

	}

	err = godotenv.Load(mflags.ConfigPath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dbUrl, err := envconfig.GetEnv("DB_URL")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	migrationsPath, err := envconfig.GetEnv("MIGRATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	migrationsTable, err := envconfig.GetEnv("MIGRATIONS_TABLE")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &MigratorConfig{
		DBUrl:           dbUrl,
		Cmd:             mflags.Cmd,
		Version:         mflags.Version,
		Steps:           mflags.Steps,
		MigrationsPath:  migrationsPath,
		MigrationsTable: migrationsTable,
	}, nil
}

type migratorFlags struct {
	Cmd        *string
	Version    *int
	Steps      *int
	ConfigPath string
}

func loadMigratorFlags() (*migratorFlags, error) {
	const op = "config.loadMigratorFlags"

	migrationCommand := flag.String("m", "up", "migration command: up, down, force, version")
	version := flag.Int("version", 0, "version for force migration")
	steps := flag.Int("steps", 0, "number of steps for up/down migration")

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configPath == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrRequiredConfigPath)
	}

	return &migratorFlags{
		Cmd:        migrationCommand,
		Version:    version,
		Steps:      steps,
		ConfigPath: configPath,
	}, nil
}

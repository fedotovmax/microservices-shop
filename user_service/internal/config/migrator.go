package config

import (
	"flag"

	"github.com/joho/godotenv"
)

type MigratorConfig struct {
	DBUrl   string
	Cmd     *string
	Version *int
	Steps   *int
}

func MustLoadMigratorConfig() *MigratorConfig {

	const op = "config.MustLoadMigratorConfig"

	mflags := loadMigratorFlags()

	ok := checkConfigPathExists(mflags.ConfigPath)

	if !ok {
		panicError(op, errConfigPathNotExists)
	}

	err := godotenv.Load(mflags.ConfigPath)

	if err != nil {
		panicError(op, err)
	}

	dbUrl, err := getEnv("DB_URL")

	if err != nil {
		panic(panicError(op, err))
	}

	return &MigratorConfig{
		DBUrl:   dbUrl,
		Cmd:     mflags.Cmd,
		Version: mflags.Version,
		Steps:   mflags.Steps,
	}
}

type migratorFlags struct {
	Cmd        *string
	Version    *int
	Steps      *int
	ConfigPath string
}

func loadMigratorFlags() *migratorFlags {
	const op = "config.loadMigratorFlags"

	migrationCommand := flag.String("m", "up", "migration command: up, down, force, version")
	version := flag.Int("version", 0, "version for force migration")
	steps := flag.Int("steps", 0, "number of steps for up/down migration")

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	if configPath == "" {
		panicError(op, errRequiredConfigPath)
	}

	return &migratorFlags{
		Cmd:        migrationCommand,
		Version:    version,
		Steps:      steps,
		ConfigPath: configPath,
	}
}

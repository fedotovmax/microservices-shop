package main

import (
	"flag"

	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	migrationCommand := flag.String("m", "up", "migration command: up, down, force, version")
	version := flag.Int("version", 0, "version for force migration")
	steps := flag.Int("steps", 0, "number of steps for up/down migration")
	flag.Parse()

	postgresUrl := os.Getenv("DB_URL")
	if postgresUrl == "" {
		slog.Error("postgres db url not provided!")
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://migrations",
		postgresUrl+"?sslmode=disable&x-migrations-table=migrations")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer m.Close()

	switch *migrationCommand {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "force":
		err = m.Force(*version)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("current migration version", "version", version, "dirty", dirty)
		return
	default:
		slog.Error("unknown migration command", "command", *migrationCommand)
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("migration completed successfully", "command", *migrationCommand)
}

package main

import (
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/user_service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	cfg := config.MustLoadMigratorConfig()

	m, err := migrate.New(
		"file://migrations",
		cfg.DBUrl+"?sslmode=disable&x-migrations-table=migrations")
	if err != nil {
		panic(err.Error())
	}
	defer m.Close()

	switch *cfg.Cmd {
	case "up":
		if *cfg.Steps > 0 {
			err = m.Steps(*cfg.Steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *cfg.Steps > 0 {
			err = m.Steps(-*cfg.Steps)
		} else {
			err = m.Down()
		}
	case "force":
		err = m.Force(*cfg.Version)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			panic(err)

		}
		slog.Info("current migration version", "version", version, "dirty", dirty)
		return
	default:
		panic(fmt.Sprintf("unknown migration command, command: %s", *cfg.Cmd))
	}

	if err != nil && err != migrate.ErrNoChange {
		panic(err.Error())
	}

	slog.Info("migration completed successfully", "command", *cfg.Cmd)
}

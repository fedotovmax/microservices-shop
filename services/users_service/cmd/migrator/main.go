package main

import (
	"fmt"
	"os"
	"path"

	"github.com/fedotovmax/microservices-shop/users_service/internal/config"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	log := logger.NewDevelopmentHandler()

	cfg, err := config.LoadMigratorConfig()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	migrationsPath := "file://" + path.Join(cfg.MigrationsPath)

	m, err := migrate.New(
		migrationsPath,
		cfg.DBUrl+"?sslmode=disable&x-migrations-table=migrations")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
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
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Info("current migration version", "version", version, "dirty", dirty)
		return
	default:
		log.Error(fmt.Sprintf("unknown migration command, command: %s", *cfg.Cmd))
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("migration completed successfully", "command", *cfg.Cmd)
}

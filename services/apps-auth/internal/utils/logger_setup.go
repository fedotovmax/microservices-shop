package utils

import (
	"log/slog"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/keys"
	"github.com/fedotovmax/microservices-shop/apps-auth/pkg/logger"
)

func SetupLooger(env string) (*slog.Logger, error) {
	switch env {
	case keys.Development:
		return logger.NewDevelopmentHandler(slog.LevelInfo), nil
	case keys.Production:
		return logger.NewProductionHandler(slog.LevelWarn), nil
	default:
		return nil, envconfig.ErrInvalidAppEnv
	}
}

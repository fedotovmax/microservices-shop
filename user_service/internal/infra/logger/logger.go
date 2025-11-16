package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/config"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func MustNewLogger(env string) *slog.Logger {

	const op = "logger.MustNewLogger"

	switch env {
	case config.Development:
		return newDevelopmentHandler()
	case config.Production:
		return newProductionHandler()
	default:
		panic(fmt.Sprintf("%s: unsopported app env for logger", op))
	}
}

func newDevelopmentHandler() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // В dev показываем больше логов
		AddSource: true,
	})
	return slog.New(handler)
}

func newProductionHandler() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("2006-01-02T15:04:05.000Z07:00"))
				}
			}
			return a
		},
	})

	return slog.New(handler)
}

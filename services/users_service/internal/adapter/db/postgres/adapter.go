package postgres

import (
	"log/slog"

	"github.com/fedotovmax/pgxtx"
)

type postgresAdapter struct {
	ex  pgxtx.Extractor
	log *slog.Logger
}

func NewPostgresAdapter(ex pgxtx.Extractor, log *slog.Logger) *postgresAdapter {
	return &postgresAdapter{
		ex:  ex,
		log: log,
	}
}

package languages

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/jackc/pgx/v5"
)

const getDefaultQuery = "select code, is_default, is_active from languages where is_default = true;"

func (p *postgres) GetDefault(ctx context.Context) (*domain.Language, error) {
	const op = "adapters.db.postgres.languages.GetDefault"

	tx := p.ex.ExtractTx(ctx)

	row := tx.QueryRow(ctx, getDefaultQuery)

	l := &domain.Language{}

	err := row.Scan(&l.Code, &l.IsDefault, &l.IsActive)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return l, nil

}

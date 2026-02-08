package languagespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
)

const getAllQuery = "select code, is_default, is_active from languages;"

func (p *postgres) GetAll(ctx context.Context) ([]domain.Language, error) {
	const op = "adapters.db.postgres.languages.GetAll"

	tx := p.ex.ExtractTx(ctx)

	rows, err := tx.Query(ctx, getAllQuery)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	defer rows.Close()

	var langs []domain.Language

	for rows.Next() {

		l := domain.Language{}

		err := rows.Scan(&l.Code, &l.IsDefault, &l.IsActive)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
		}

		langs = append(langs, l)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return langs, nil

}

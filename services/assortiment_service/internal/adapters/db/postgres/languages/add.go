package languages

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
)

const addQuery = "insert into languages (code, is_default, is_active) values ($1,$2,$3);"

func (p *postgres) Add(ctx context.Context, code string, isDefault, isActive bool) error {
	const op = "adapters.db.postgres.languages.Add"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, addQuery, code, isDefault, isActive)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}

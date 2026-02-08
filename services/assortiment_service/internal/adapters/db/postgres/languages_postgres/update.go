package languagespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
)

const updateQuery = "update languages set is_default = $1, is_active = $2 where code = $3;"

func (p *postgres) Update(ctx context.Context, code string, isDefault, isActive bool) error {

	const op = "adapters.db.postgres.languages.Update"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateQuery, isDefault, isActive, code)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}

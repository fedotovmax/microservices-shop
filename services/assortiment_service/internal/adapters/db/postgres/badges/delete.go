package badges

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
)

const deleteQuery = "delete from badges where code = $1;"

func (p *postgres) Delete(ctx context.Context, code string) error {
	const op = "adapters.db.postgres.badges.delete"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, deleteQuery, code)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}

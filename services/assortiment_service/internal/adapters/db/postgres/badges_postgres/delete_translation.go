package badgespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
)

const deleteTranslationQuery = `delete from badge_translations where id = $1;`

func (p *postgres) DeleteTranslation(ctx context.Context, id string) error {
	const op = "adapters.db.postgres.badges.DeleteTranslation"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, deleteTranslationQuery, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}

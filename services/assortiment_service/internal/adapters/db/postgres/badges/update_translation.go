package badges

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const updateTranslationQuery = `
update badge_translations
set title = $1
where id = $2;`

func (p *postgres) UpdateTranslation(ctx context.Context, in *inputs.UpdateBadgeTranslate) error {
	const op = "adapters.db.postgres.badges.UpdateTranslation"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateTranslationQuery, in.Title, in.ID)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}

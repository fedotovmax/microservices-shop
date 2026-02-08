package badgespostgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

const saveQuery = `
insert into badges
(code, starts_at, ends_at, color, priority)
values ($1, $2, $3, $4, $5);`

func (p *postgres) Save(ctx context.Context, in *inputs.SaveBadge) error {

	const op = "adapters.db.postgres.badges.Save"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, saveQuery, in.Code, in.StartsAt, in.EndsAt, in.Color, in.Priority)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}

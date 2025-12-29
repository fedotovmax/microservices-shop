package postgres

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

const setEventStatusDoneQuery = "update events set status = $1 where id = $2;"

func (p *postgresAdapter) SetEventStatusDone(ctx context.Context, id string) error {
	const op = "adapter.db.postgres.SetEventStatusDone"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setEventStatusDoneQuery, domain.EventStatusDone, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

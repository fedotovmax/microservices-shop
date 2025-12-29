package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
)

const setEventsReservedToByIDsQuery = "update events set reserved_to = $1 where id = ANY ($2);"

func (p *postgresAdapter) SetEventsReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error {

	const op = "adapter.db.postgres.SetEventsReservedToByIDs"

	reservedTo := time.Now().Add(dur)

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, setEventsReservedToByIDsQuery, reservedTo, ids)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}

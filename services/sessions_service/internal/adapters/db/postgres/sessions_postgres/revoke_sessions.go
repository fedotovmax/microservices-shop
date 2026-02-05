package sessionspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
)

const revokeSessionQuery = "update sessions set revoked_at = $1 where id = any($2) and revoked_at is null;"

func (p *postgres) Revoke(ctx context.Context, sids []string) error {

	const op = "adapter.db.postgres.sessions.Revoke"

	if len(sids) == 0 {
		return nil
	}

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, revokeSessionQuery, time.Now().UTC(), sids)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}

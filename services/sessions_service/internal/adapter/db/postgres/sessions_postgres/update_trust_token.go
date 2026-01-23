package sessionspostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const updateTrustTokenQuery = `update trust_tokens
set last_used_at = $1,
    expires_at = $2
where token_hash = $3
  and revoked_at is null
  and expires_at > $1;`

func (p *postgres) UpdateTrustToken(ctx context.Context, in *inputs.CreateTrustTokenInput) error {
	const op = "adapter.db.postgres.UpdateTrustToken"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, updateTrustTokenQuery, time.Now().UTC(), in.ExpiresAt, in.TokenHash)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

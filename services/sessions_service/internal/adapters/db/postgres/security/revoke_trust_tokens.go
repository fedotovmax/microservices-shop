package security

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
)

const revokeTrustTokensQuery = "update trust_tokens set revoked_at = $1 where token_hash = any($2) and revoked_at is null;"

func (p *postgres) RevokeTrustTokens(ctx context.Context, hashes []string) error {

	const op = "adapter.db.postgres.RevokeTrustTokens"

	if len(hashes) == 0 {
		return nil
	}

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, revokeTrustTokensQuery, time.Now().UTC(), hashes)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}

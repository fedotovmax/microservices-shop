package securitypostgres

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

const createTrustToken = `
insert into trust_tokens (token_hash, uid, last_used_at, expires_at) values ($1, $2, $3, $4);`

func (p *postgres) CreateTrustToken(ctx context.Context, in *inputs.CreateTrustTokenInput) error {

	const op = "adapter.db.postgres.CreateTrustToken"

	tx := p.ex.ExtractTx(ctx)

	_, err := tx.Exec(ctx, createTrustToken, in.TokenHash, in.UID, time.Now().UTC(), in.ExpiresAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil

}

package security

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
)

const findUserTrustTokensQuery = "select token_hash, uid, last_used_at, expires_at, revoked_at from trust_tokens where uid = $1 and revoked_at is null and expires_at > $2 order by last_used_at desc;"

func (p *postgres) FindUserTrustTokens(ctx context.Context, uid string) ([]*domain.DeviceTrustToken, error) {
	const op = "adapter.db.postgres.FindUserTrustTokens"

	tx := p.ex.ExtractTx(ctx)

	rows, err := tx.Query(ctx, findUserTrustTokensQuery, uid, time.Now().UTC())

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	defer rows.Close()

	var tokens []*domain.DeviceTrustToken

	for rows.Next() {
		t := &domain.DeviceTrustToken{}

		err := rows.Scan(
			&t.TokenHash, &t.UID,
			&t.LastUsed, &t.ExpiresAt,
			&t.RevokedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
		}

		tokens = append(tokens, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return tokens, nil

}
